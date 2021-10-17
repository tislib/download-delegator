package app

import (
	"bytes"
	"compress/gzip"
	"download-delegator/model"
	"download-delegator/service"
	"encoding/json"
	"log"
	"net"
	"net/http"
	pprof "net/http/pprof"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type App struct {
	srv          *http.Server
	Async        bool
	config       model.Config
	Addr         string
	pprofHandler http.Handler
}

func (app *App) Init(config model.Config) {
	app.config = config

	app.Addr = config.Listen.Addr
}

func (app *App) Run() {
	app.srv = &http.Server{Addr: app.config.Listen.Addr, Handler: app}

	service.DownloaderServiceInstance.ProxyFile = app.config.Proxy.File

	service.DownloaderServiceInstance.ConfigureSanitizer()

	app.pprofHandler = pprof.Handler("pprof")

	app.startListening()
}

func (app *App) ListenAndServeAsync() {
	srv := app.srv
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	app.Addr = ln.Addr().String()

	if app.Async {
		go func() {
			log.Fatal(srv.Serve(ln))
		}()
	} else {
		log.Fatal(srv.Serve(ln))
	}
}

func (app *App) startListening() {
	log.Printf("Serving on " + app.Addr)
	app.ListenAndServeAsync()
}

func (app *App) GetAddr() string {
	return app.Addr
}

func (app *App) Close() error {
	return app.srv.Close()
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Method + " " + r.URL.Path

	log.Print("download: ", r.RequestURI, " ", r.RemoteAddr)

	if strings.HasPrefix(action, "GET /pprof") {
		app.pprofHandler.ServeHTTP(w, r)
	}

	if strings.HasPrefix(action, "POST /pprof") {
		app.pprofHandler.ServeHTTP(w, r)
	}

	switch action {
	case "GET /test":
		app.test(w, r)
		break
	case "GET /get":
		status := app.get(w, r, false)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /get":
		status := app.get(w, r, true)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /bulk":
		status := app.bulk(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	}
}

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(100 * time.Second)
	w.WriteHeader(404)

	w.Write([]byte("hello world"))
}

func (app *App) bulk(w http.ResponseWriter, r *http.Request) int {
	defer r.Body.Close()

	timeCalc := new(TimeCalc)
	timeCalc.Init("bulk-download")

	var config model.BulkDownloadConfig

	err := json.NewDecoder(r.Body).Decode(&config)

	log.Print(config)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		writeDownloadError(w, err, &model.DownloadError{
			ErrorState:   model.RequestBodyNotValid,
			ErrorText:    err.Error(),
			ClientStatus: 0,
		})

		return 400
	}

	wg := new(sync.WaitGroup)

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	var maxConcurrencySemaphore = make(chan int, config.MaxConcurrency)

	resultChan := make(chan model.DownloadResponse, config.MaxConcurrency*2)

	go func() {
		var counter int32
		for indexX, itemX := range config.Url {
			wg.Add(1)
			maxConcurrencySemaphore <- 1

			go func(index int, item string) {
				atomic.AddInt32(&counter, 1)

				defer func() {
					atomic.AddInt32(&counter, -1)
					log.Print("concurrency level: ", strconv.Itoa(int(counter)))

					if err := recover(); err != nil {
						log.Print(err)
					}
					wg.Done()
					<-maxConcurrencySemaphore
				}()

				select {
				case <-r.Context().Done():
					return
				default:

				}

				log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"begin bulk download index/url: ", index, item)

				downloadConfig := model.DownloadConfig{
					Url:      item,
					Compress: false,
					Sanitize: config.Sanitize,
					Proxy:    config.Proxy,
					Timeout:  config.Timeout,
				}

				var buf bytes.Buffer

				beginTime := time.Now()

				var resItem model.DownloadResponse

				for i := 0; i < config.RetryCount; i++ {
					statusCode, downloadErr, err := service.DownloaderServiceInstance.Get(&buf, r.Context(), downloadConfig)

					duration := time.Now().Sub(beginTime)

					localResItem := model.DownloadResponse{
						Url:           item,
						Index:         index,
						StatusCode:    statusCode,
						Content:       buf.String(),
						DownloadError: downloadErr,
						Duration:      duration,
						Retried:       i,
						DurationMS:    int(duration / time.Millisecond),
					}

					resItem = localResItem

					if err != nil {
						log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"end bulk download index/url: ", index, item, statusCode, len(resItem.Content), int(duration/time.Millisecond))
						break
					}
				}

				timeCalc.Step()
				log.Print("sending to chan", counter)
				resultChan <- resItem
				log.Print("sent to chan", counter)

				if err != nil {
					log.Print(err)
				}
			}(indexX, itemX)
		}

		wg.Wait()
		close(resultChan)
	}()

	bodyWriter := w

	if config.Compress {
		w.Header().Set("Content-Encoding", "application/gzip")

		gzipWriter := gzip.NewWriter(w)

		defer func() {
			err := gzipWriter.Close()

			if err != nil {
				log.Print(err)
			}
		}()

		w = bodyWriter
	}

	if config.OutputForm == "" || config.OutputForm == model.JsonOutput {
		w.Header().Set("Content-Type", "application/json")

		bodyWriter.Write([]byte("["))
		isFirst := true

		for item := range resultChan {
			log.Print("begin write for", item.Url)
			data, err := json.Marshal(item)

			if err != nil {
				log.Print(err)
			}

			if !isFirst {
				bodyWriter.Write([]byte(",\n"))
			}

			_, err = w.Write(data)

			if err != nil {
				log.Print(err)
			}

			isFirst = false
			log.Print("end write", item.Url)
		}

		bodyWriter.Write([]byte("]"))

	}

	return 200
}

func (app *App) get(w http.ResponseWriter, r *http.Request, useBody bool) int {
	defer r.Body.Close()

	var config model.DownloadConfig
	if !useBody {
		query, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)

			writeDownloadError(w, err, &model.DownloadError{
				ErrorState:   model.UrlNotValid,
				ErrorText:    err.Error(),
				ClientStatus: 0,
			})

			return 400
		}

		config = app.parseConfig(query, err)
	} else {
		err := json.NewDecoder(r.Body).Decode(&config)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)

			writeDownloadError(w, err, &model.DownloadError{
				ErrorState:   model.RequestBodyNotValid,
				ErrorText:    err.Error(),
				ClientStatus: 0,
			})

			return 400
		}
	}

	statusCode, downloadErr, err := service.DownloaderServiceInstance.Get(w, r.Context(), config)

	if downloadErr != nil || err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		writeDownloadError(w, err, downloadErr)

		return 400
	}

	if config.Compress {
		w.Header().Set("Content-Encoding", "gzip")
	}

	return statusCode
}

func (app *App) parseConfig(query url.Values, err error) model.DownloadConfig {
	urlParam := query.Get("url")

	timeout, err := strconv.Atoi(query.Get("timeout"))
	if err != nil {
		timeout = 5000
	}

	config := model.DownloadConfig{
		Url:   urlParam,
		Proxy: query.Get("proxy") == "true",
		Timeout: model.TimeoutConfig{
			RequestTimeout: time.Duration(timeout) * time.Millisecond,
		},
		Compress: query.Get("compress") == "true",
		Sanitize: model.SanitizeConfig{
			CleanMinimal:  query.Get("cleanMinimal") == "true",
			CleanMinimal2: query.Get("cleanMinimal2") == "true",
		},
	}
	return config
}

func writeDownloadError(w http.ResponseWriter, err error, downloadError *model.DownloadError) {
	log.Print(err)

	bytes, err := json.Marshal(downloadError)

	if err != nil {
		log.Print(err)
	}

	_, err = w.Write(bytes)

	if err != nil {
		log.Print(err)
	}
}

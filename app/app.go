package app

import (
	"bytes"
	"download-delegator/model"
	"download-delegator/service"
	"encoding/json"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type App struct {
	srv    *http.Server
	Async  bool
	config model.Config
	Addr   string
}

func (app *App) Init(config model.Config) {
	app.config = config

	app.Addr = config.Listen.Addr
}

func (app *App) Run() {
	app.srv = &http.Server{Addr: app.config.Listen.Addr, Handler: app}

	service.DownloaderServiceInstance.ProxyFile = app.config.Proxy.File

	service.DownloaderServiceInstance.ConfigureSanitizer()

	app.startListening()
}

func (app *App) ListenAndServeAsync() {
	srv := app.srv
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	app.Addr = ln.Addr().String()

	if app.Async {
		go func() {
			log.Fatal(srv.ServeTLS(ln, app.config.Tls.Cert, app.config.Tls.Key))
		}()
	} else {
		log.Fatal(srv.ServeTLS(ln, app.config.Tls.Cert, app.config.Tls.Key))
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

	var config model.BulkDownloadConfig

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

	var result []model.DownloadResponse

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	var maxConcurrencySemaphore = make(chan int, config.MaxConcurrency)

	for indexX, itemX := range config.Url {
		wg.Add(1)
		maxConcurrencySemaphore <- 1

		go func(index int, item string) {
			defer func() {
				wg.Done()
				<-maxConcurrencySemaphore
			}()

			select {
			case <-r.Context().Done():
				return
			default:

			}

			log.Print("begin bulk download index/url: ", index, item)

			downloadConfig := model.DownloadConfig{
				Url:      item,
				Compress: false,
				Sanitize: config.Sanitize,
				Proxy:    config.Proxy,
				Timeout:  config.Timeout,
			}

			var buf bytes.Buffer

			beginTime := time.Now()

			statusCode, downloadErr, err := service.DownloaderServiceInstance.Get(&buf, r.Context(), downloadConfig)

			duration := time.Now().Sub(beginTime)

			resItem := model.DownloadResponse{
				Url:           item,
				Index:         index,
				StatusCode:    statusCode,
				Content:       buf.String(),
				DownloadError: downloadErr,
				Duration:      duration,
				DurationMS:    int(duration / time.Millisecond),
			}

			log.Print("end bulk download index/url: ", index, item, statusCode, len(resItem.Content), int(duration/time.Millisecond))

			mutex.Lock()
			result = append(result, resItem)
			mutex.Unlock()

			if err != nil {
				log.Print(err)
			}
		}(indexX, itemX)
	}

	wg.Wait()

	if config.OutputForm == "" || config.OutputForm == model.JsonOutput {
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(result)

		if err != nil {
			log.Print(err)
		}

		_, err = w.Write(data)

		if err != nil {
			log.Print(err)
		}
	}

	if config.Compress {
		w.Header().Set("Content-Encoding", "gzip")
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
		Url:      urlParam,
		Proxy:    query.Get("proxy") == "true",
		Timeout:  time.Duration(timeout) * time.Millisecond,
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

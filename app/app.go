package app

import (
	"compress/gzip"
	"download-delegator/model"
	error2 "download-delegator/model/errors"
	"download-delegator/service"
	"encoding/json"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	"io"
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
	Version      string
}

func (app *App) Init(config model.Config) {
	app.config = config

	app.Addr = config.Listen.Addr
}

func (app *App) Run() {
	app.srv = &http.Server{Addr: app.config.Listen.Addr, Handler: app}

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
	case "GET /version":
		status := app.version(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "GET /get":
		status := app.get(w, r, false)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /get":
		status := app.get(w, r, true)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /bulk":
		status := app.bulkDownload(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /whois":
		status := app.bulkWhois(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	}
}

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(100 * time.Second)
	w.WriteHeader(404)

	w.Write([]byte("hello world"))
}

func (app *App) bulkDownload(w http.ResponseWriter, r *http.Request) int {
	defer r.Body.Close()

	timeCalc := new(TimeCalc)
	timeCalc.Init("bulk-download")

	var config model.BulkDownloadConfig

	err := json.NewDecoder(r.Body).Decode(&config)

	log.Print(config)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		log.Print(err)

		writeDownloadError(w, err, error2.RequestBodyNotValid)

		return 400
	}

	wg := new(sync.WaitGroup)

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	var maxConcurrencySemaphore = make(chan int, config.MaxConcurrency)

	resultChan := make(chan model.DownloadResponse, config.MaxConcurrency*2)

	downloaderService := new(service.DownloaderService)
	downloaderService.InitTransformers(config.Transform)

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
					Url:       item,
					Proxy:     config.Proxy,
					Transform: config.Transform,
					Timeout:   config.Timeout,
				}

				var resItem model.DownloadResponse

				for i := 0; i < config.RetryCount; i++ {
					resItem = downloaderService.Get(r.Context(), downloadConfig)

					resItem.Index = index
					resItem.Retried = i

					if err != nil {
						log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"end bulk download index/url: ", index, item, resItem.StatusCode, len(resItem.Content), resItem.DurationMS)
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

	var bodyWriter io.Writer = w

	if config.Compression.IsCompressionEnabled() {
		writer := app.compress(w, config.Compression)
		bodyWriter = writer

		defer func() {
			err := writer.Close()

			if err != nil {
				log.Print(err)
			}
		}()
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

			_, err = bodyWriter.Write(data)

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

func (app *App) compress(w http.ResponseWriter, compression model.Compression) io.WriteCloser {
	w.Header().Set("Content-Encoding", "application/gzip")

	var writer io.WriteCloser
	var err error

	switch compression.Algo {
	case model.Gzip:
		writer, err = gzip.NewWriterLevel(w, compression.Level)
		break
	case model.Bzip2:
		writer, err = bzip2.NewWriter(w, &bzip2.WriterConfig{Level: compression.Level})
		break
	case model.Zstd:
		writer, err = zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(compression.Level)))
		break
	}

	if err != nil {
		log.Print(err)
	}

	return writer
}
func (app *App) bulkWhois(w http.ResponseWriter, r *http.Request) int {
	defer r.Body.Close()

	timeCalc := new(TimeCalc)
	timeCalc.Init("bulk-whois")

	var config model.BulkWhoisConfig

	err := json.NewDecoder(r.Body).Decode(&config)

	log.Print(config)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		log.Println(err)

		writeDownloadError(w, err, error2.RequestBodyNotValid)

		return 400
	}

	wg := new(sync.WaitGroup)

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	var maxConcurrencySemaphore = make(chan int, config.MaxConcurrency)

	resultChan := make(chan model.WhoisResponse, config.MaxConcurrency*2)

	go func() {
		var counter int32
		for indexX, itemX := range config.Domains {
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

				log.Print("["+(strconv.Itoa(len(config.Domains)))+"/"+strconv.Itoa(index)+"]"+"begin bulk download index/url: ", index, item)

				var resItem model.WhoisResponse

				for i := 0; i < config.RetryCount; i++ {
					resItem = service.WhoisServiceInstance.Get(item, config.Timeout)

					if resItem.Error == error2.NoError {
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

	var bodyWriter io.Writer = w

	if config.Compression.IsCompressionEnabled() {
		writer := app.compress(w, config.Compression)
		bodyWriter = writer

		defer func() {
			err := writer.Close()

			if err != nil {
				log.Print(err)
			}
		}()
	}

	if config.OutputForm == "" || config.OutputForm == model.JsonOutput {
		w.Header().Set("Content-Type", "application/json")

		bodyWriter.Write([]byte("["))
		isFirst := true

		for item := range resultChan {
			log.Print("begin write for", item.Domain)
			data, err := json.Marshal(item)

			if err != nil {
				log.Print(err)
			}

			if !isFirst {
				bodyWriter.Write([]byte(",\n"))
			}

			_, err = bodyWriter.Write(data)

			if err != nil {
				log.Print(err)
			}

			isFirst = false
			log.Print("end write", item.Domain)
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

			log.Println(err)

			writeDownloadError(w, err, error2.UrlNotValid)

			return 400
		}

		config = app.parseConfig(query, err)
	} else {
		err := json.NewDecoder(r.Body).Decode(&config)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)

			log.Print(err)

			writeDownloadError(w, err, error2.RequestBodyNotValid)

			return 400
		}
	}

	downloaderService := new(service.DownloaderService)
	downloaderService.InitTransformers(config.Transform)

	downloadResponse := downloaderService.Get(r.Context(), config)

	if downloadResponse.Error != error2.NoError {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		writeDownloadError(w, nil, downloadResponse.Error)

		return 400
	}

	w.Write([]byte(downloadResponse.Content))

	return downloadResponse.StatusCode
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
	}
	return config
}

func (app *App) version(w http.ResponseWriter, r *http.Request) int {
	w.Write([]byte(app.Version))

	return 200
}

func writeDownloadError(w http.ResponseWriter, err error, error error2.State) {
	_, err = w.Write([]byte(error))

	if err != nil {
		log.Print(err)
	}
}

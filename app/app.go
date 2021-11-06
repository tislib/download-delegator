package app

import (
	"compress/gzip"
	model2 "download-delegator/core/model"
	error2 "download-delegator/core/model/errors"
	"download-delegator/core/service"
	service2 "download-delegator/lib/impl"
	"encoding/json"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	log "github.com/sirupsen/logrus"
	"io"
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
	config       model2.Config
	addr         string
	pprofHandler http.Handler
	Version      string
	appService   service.AppService
}

func (app *App) Init(config model2.Config) {
	app.config = config

	app.addr = config.Listen.Addr

	app.appService = new(service2.AppServiceImpl)
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
	app.addr = ln.Addr().String()

	if app.Async {
		go func() {
			log.Fatal(srv.Serve(ln))
		}()
	} else {
		log.Fatal(srv.Serve(ln))
	}
}

func (app *App) startListening() {
	log.Printf("Serving on " + app.addr)
	app.ListenAndServeAsync()
}

func (app *App) GetAddr() string {
	return app.addr
}

func (app *App) Close() error {
	return app.srv.Close()
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Method + " " + r.URL.Path

	log.Debug("download: ", r.RequestURI, " ", r.RemoteAddr)

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
		log.Debug("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "GET /get":
		status := app.get(w, r, false)
		log.Debug("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /get":
		status := app.get(w, r, true)
		log.Debug("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /bulk":
		status := app.bulkDownload(w, r)
		log.Debug("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "POST /whois":
		status := app.bulkWhois(w, r)
		log.Debug("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
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

	var config model2.BulkDownloadConfig

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

	resultChan := make(chan model2.DownloadResponse, config.MaxConcurrency*2)

	downloaderService := new(service2.DownloaderService)
	downloaderService.ConfigureTransformers(config.Transform)
	downloaderService.ConfigureTimeout(config.Timeout)
	downloaderService.EnableProxy(config.Proxy)

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

				var resItem model2.DownloadResponse

				for i := 0; i < config.RetryCount; i++ {
					resItem = downloaderService.Get(r.Context(), item)

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

	if config.OutputForm == "" || config.OutputForm == model2.JsonOutput {
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

func (app *App) compress(w http.ResponseWriter, compression model2.Compression) io.WriteCloser {
	w.Header().Set("Content-Encoding", "application/gzip")

	var writer io.WriteCloser
	var err error

	switch compression.Algo {
	case model2.Gzip:
		writer, err = gzip.NewWriterLevel(w, compression.Level)
		break
	case model2.Bzip2:
		writer, err = bzip2.NewWriter(w, &bzip2.WriterConfig{Level: compression.Level})
		break
	case model2.Zstd:
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

	var config model2.BulkWhoisConfig

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

	resultChan := make(chan model2.WhoisResponse, config.MaxConcurrency*2)

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

				var resItem model2.WhoisResponse

				for i := 0; i < config.RetryCount; i++ {
					resItem = service2.WhoisServiceInstance.Get(item, config.Timeout)

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

	if config.OutputForm == "" || config.OutputForm == model2.JsonOutput {
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

	jsonEncoder := json.NewEncoder(w)

	sendError := func(errorState error2.State) int {
		downloadResponse := new(model2.DownloadResponse)
		downloadResponse.Error = errorState

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		_ = jsonEncoder.Encode(downloadResponse)

		return 400
	}

	var config model2.DownloadConfig
	if !useBody {
		query, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			log.Warn(err)

			return sendError(error2.UrlNotValid)
		}

		config = app.parseConfig(query, err)
	} else {
		err := json.NewDecoder(r.Body).Decode(&config)

		if err != nil {
			log.Warn(err)

			return sendError(error2.RequestBodyNotValid)
		}
	}

	downloadResponse, err := app.appService.Get(r.Context(), config)

	if err != nil {
		log.Warn(err)
	}

	_ = jsonEncoder.Encode(downloadResponse)

	return downloadResponse.StatusCode
}

func (app *App) parseConfig(query url.Values, err error) model2.DownloadConfig {
	urlParam := query.Get("url")

	timeout, err := strconv.Atoi(query.Get("timeout"))
	if err != nil {
		timeout = 5000
	}

	config := model2.DownloadConfig{
		Url:   urlParam,
		Proxy: query.Get("proxy") == "true",
		Timeout: model2.TimeoutConfig{
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

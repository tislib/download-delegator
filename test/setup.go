package test

import (
	"crypto/tls"
	app2 "download-delegator/app"
	client2 "download-delegator/client"
	"download-delegator/core/model"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"time"
)

func runServer() *app2.App {
	app := new(app2.App)

	app.Async = true

	app.Init(model.Config{
		Concurrency: model.ConcurrencyConfig{
			MaxConcurrency: 100,
		},
		Listen: model.ListenConfig{
			Addr: ":0",
		},
		Proxy: model.ProxyConfig{},
	})

	app.Run()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return app
}

type lightServer struct {
	Port int
}

func (receiver *lightServer) run() {
	srv := &http.Server{Handler: receiver}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Panic(err)
	}

	receiver.Port = ln.Addr().(*net.TCPAddr).Port

	go func() {
		log.Fatal(srv.Serve(ln))
	}()
}

func (receiver lightServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Method + " " + r.URL.Path

	switch action {
	case "GET /get-raw":
		func() {
			data := r.URL.Query().Get("data")
			repeatValStr := r.URL.Query().Get("repeat")
			repeatVal := 1
			if repeatValStr != "" {
				repeatVal, _ = strconv.Atoi(repeatValStr)
			}

			for i := 0; i < repeatVal; i++ {
				_, _ = w.Write([]byte(data))
			}
		}()
		break
	case "GET /timeout":
		func() {
			timeValue, err := strconv.Atoi(r.URL.Query().Get("val"))

			if err != nil {
				log.Panic(err)
			}

			time.Sleep(time.Millisecond * time.Duration(timeValue))
		}()
		break
	}
}

var app = runServer()
var Client *client2.DownloadDelegatorClient

var LightTestServer = lightServer{}

func init() {
	log.Print("test init")

	Client = new(client2.DownloadDelegatorClient)
	Client.Init(client2.InitConfig{Addr: app.GetAddr()})

	log.SetLevel(log.TraceLevel)

	LightTestServer.run()
}

package app

import (
	"bytes"
	"crypto/tls"
	"download-delegator/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	url2 "net/url"
	"testing"
)

func runServer() *App {
	app := new(App)

	app.Addr = ":0"
	app.Async = true

	app.Run()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return app
}

func TestDnsResolveProblem(t *testing.T) {
	app := runServer()

	log.Print(app.GetAddr())

	config := model.DownloadConfig{
		Url:   "non-existing-domain-123.com",
		Proxy: true,
	}

	var buf bytes.Buffer

	statusCode, err := Get(app, &buf, config)

	if err != nil {
		log.Panic(err)
	}

	log.Print("FINISH!")
	log.Print(statusCode, buf.String())
}

func Get(app *App, w io.Writer, config model.DownloadConfig) (int, error) {
	url := "http://" + app.GetAddr() + "/get"

	url += "?url=" + url2.PathEscape(config.Url)
	if !config.Proxy {
		url += "&noProxy"
	}

	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	_, err = io.Copy(w, resp.Body)

	return resp.StatusCode, err
}

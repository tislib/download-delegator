package app

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"download-delegator/model"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"time"
)

type App struct {
	Addr      string
	CertFile  string
	KeyFile   string
	ProxyFile string

	proxyList []model.ProxyConfig
	sanitizer *bluemonday.Policy
}

func (app *App) loadProxyConfig() {
	csvFile, err := os.Open(app.ProxyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, line := range csvLines {
		proxyConfig := model.ProxyConfig{
			Host:     line[0],
			Port:     line[1],
			Username: line[2],
			Password: line[3],
		}

		app.proxyList = append(app.proxyList, proxyConfig)
	}

	log.Printf("Proxies loaded: %d proxy", len(app.proxyList))
}

func (app *App) Run() {
	app.loadProxyConfig()

	app.configureSanitizer()

	app.startListening()
}

func (app *App) startListening() {
	srv := &http.Server{Addr: app.Addr, Handler: app}
	log.Printf("Serving on " + app.Addr)
	log.Fatal(srv.ListenAndServeTLS(app.CertFile, app.KeyFile))
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Method + " " + r.URL.Path

	log.Print("download: ", r.RequestURI, " ", r.RemoteAddr)

	switch action {
	case "GET /test":
		app.test(w, r)
		break
	case "GET /get":
		status := app.get(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	case "GET /get-clean":
		status := app.getClean(w, r)
		log.Print("result: ", r.RequestURI, " ", r.RemoteAddr, " ", status)
	}
}

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Write([]byte(r.RequestURI))
}

func (app *App) get(w http.ResponseWriter, r *http.Request) uint64 {
	defer r.Body.Close()
	gzw := gzip.NewWriter(w)

	query, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		log.Print("url parse error: ", err)
		return 500
	}

	urlParam := query.Get("url")

	if urlParam == "" {
		w.WriteHeader(404)
		w.Write([]byte("invalid url"))
		return 404
	}

	client := new(http.Client)
	client.Timeout = time.Second * 100

	app.configureProxy(client)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-r.Context().Done()
		cancel()
	}()

	go func() {
		<-ctx.Done()
		r.Body.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", urlParam, nil)

	if err != nil {
		log.Print("request creation error", err)
		w.WriteHeader(400)
		return 400
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Print("request execution error", err)
		w.WriteHeader(500)
		return 500
	}

	defer resp.Body.Close()
	defer gzw.Close()

	w.Header().Set("Content-Encoding", "gzip")
	_, err = io.CopyN(gzw, resp.Body, 1024*1024*1024)

	if err != nil && err != io.EOF {
		log.Print("io error", err)
		w.WriteHeader(500)
		return 500
	}

	return 200
}
func (app *App) getClean(w http.ResponseWriter, r *http.Request) uint64 {
	defer r.Body.Close()
	gzw := gzip.NewWriter(w)

	query, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		log.Print("url parse error: ", err)
		return 500
	}

	urlParam := query.Get("url")

	if urlParam == "" {
		w.WriteHeader(404)
		w.Write([]byte("invalid url"))
		return 404
	}

	client := new(http.Client)
	client.Timeout = time.Second * 100

	app.configureProxy(client)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-r.Context().Done()
		cancel()
	}()

	go func() {
		<-ctx.Done()
		r.Body.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", urlParam, nil)

	if err != nil {
		log.Print("request creation error", err)
		w.WriteHeader(400)
		return 400
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Print("request execution error", err)
		w.WriteHeader(500)
		return 500
	}

	defer resp.Body.Close()
	defer gzw.Close()

	w.Header().Set("Content-Encoding", "gzip")

	body, err := ioutil.ReadAll(resp.Body)

	//gzw.Write(body)
	//gzw.Write(gohtml.FormatBytes(body))
	gzw.Write(gohtml.FormatBytes(app.sanitizer.SanitizeBytes(gohtml.FormatBytes(body))))

	if err != nil && err != io.EOF {
		log.Print("io error", err)
		w.WriteHeader(500)
		return 500
	}

	return 200
}

func (app *App) configureProxy(client *http.Client) {
	proxyConfig := app.locateRandomProxy()

	if proxyConfig == nil {
		return
	}

	proxyUrl, err := url.Parse("http://" + proxyConfig.Host + ":" + proxyConfig.Port)

	if err != nil {
		log.Print(err)
		return
	}

	auth := proxyConfig.Username + ":" + proxyConfig.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
		ProxyConnectHeader: map[string][]string{
			"Proxy-Authorization": append([]string{}, basicAuth),
		},
	}
}

func (app *App) locateRandomProxy() *model.ProxyConfig {
	if len(app.proxyList) > 0 {
		randomIndex := rand.Intn(len(app.proxyList))
		return &app.proxyList[randomIndex]
	}

	return nil
}

func (app *App) configureSanitizer() {
	app.sanitizer = bluemonday.NewPolicy()

	// Require URLs to be parseable by net/url.Parse and either:
	//   mailto: http:// or https://
	app.sanitizer.AllowStandardURLs()

	// We only allow <p> and <a href="">
	app.sanitizer.AllowAttrs("href").OnElements("a")
	app.sanitizer.AllowAttrs("name", "content", "property").OnElements("meta")
	app.sanitizer.AllowElements("meta", "a", "html", "head", "body", "title")
	app.sanitizer.AllowLists()
	app.sanitizer.AllowTables()
}

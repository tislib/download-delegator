package app

import (
	"compress/gzip"
	"crypto/tls"
	"download-delegator/model"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
)

type App struct {
	Addr      string
	CertFile  string
	KeyFile   string
	ProxyFile string

	proxyList []model.ProxyConfig
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

	app.startListening()
}

func (app *App) startListening() {
	srv := &http.Server{Addr: app.Addr, Handler: app}
	log.Printf("Serving on " + app.Addr)
	log.Fatal(srv.ListenAndServeTLS(app.CertFile, app.KeyFile))
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Method + " " + r.URL.Path

	switch action {
	case "GET /test":
		app.test(w, r)
		break
	case "GET /get":
		app.get(w, r)
	}
}

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Write([]byte(r.RequestURI))
}

func (app *App) get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	gzw := gzip.NewWriter(w)

	query, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		w.Write([]byte(err.Error()))
		log.Print(err)
		return
	}

	urlParam := query.Get("url")

	if urlParam == "" {
		w.WriteHeader(404)
		w.Write([]byte("invalid url"))
		return
	}

	client := new(http.Client)

	app.configureProxy(client)

	req, err := http.NewRequest("GET", urlParam, nil)

	if err != nil {
		log.Print(err)
		w.WriteHeader(400)
		return
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}

	defer resp.Body.Close()
	defer gzw.Close()

	w.Header().Set("Content-Encoding", "gzip")
	_, err = io.CopyN(gzw, resp.Body, 1024*1024*1024)

	if err != nil && err != io.EOF {
		log.Print(err)
		w.WriteHeader(500)
		return
	}

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

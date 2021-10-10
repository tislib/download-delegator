package service

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"download-delegator/model"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

type downloaderService struct {
	ProxyFile string
	proxyList []model.ProxyItemConfig
	sanitizer *bluemonday.Policy
}

func (s *downloaderService) ConfigureSanitizer() {
	s.sanitizer = bluemonday.NewPolicy()

	// Require URLs to be parseable by net/url.Parse and either:
	//   mailto: http:// or https://
	s.sanitizer.AllowStandardURLs()

	// We only allow <p> and <a href="">
	s.sanitizer.AllowAttrs("href").OnElements("a")
	s.sanitizer.AllowAttrs("name", "content", "property").OnElements("meta")
	s.sanitizer.AllowElements("meta", "a", "html", "head", "body", "title")
	s.sanitizer.AllowLists()
	s.sanitizer.AllowTables()
}

func (s *downloaderService) loadProxyConfig() {
	csvFile, err := os.Open(s.ProxyFile)
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
		ProxyItemConfig := model.ProxyItemConfig{
			Host:     line[0],
			Port:     line[1],
			Username: line[2],
			Password: line[3],
		}

		s.proxyList = append(s.proxyList, ProxyItemConfig)
	}

	log.Printf("Proxies loaded: %d proxy", len(s.proxyList))
}

func (s *downloaderService) locateRandomProxy() *model.ProxyItemConfig {
	if len(s.proxyList) > 0 {
		randomIndex := rand.Intn(len(s.proxyList))
		return &s.proxyList[randomIndex]
	}

	return nil
}

func (s *downloaderService) configureProxy(client *http.Client) {
	ProxyItemConfig := s.locateRandomProxy()

	if ProxyItemConfig == nil {
		return
	}

	proxyUrl, err := url.Parse("http://" + ProxyItemConfig.Host + ":" + ProxyItemConfig.Port)

	if err != nil {
		log.Print(err)
		return
	}

	auth := ProxyItemConfig.Username + ":" + ProxyItemConfig.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	client.Transport = s.configureTransport(&http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		ProxyConnectHeader: map[string][]string{
			"Proxy-Authorization": append([]string{}, basicAuth),
		},
	})
}

func (s *downloaderService) configureTransport(transport *http.Transport) *http.Transport {
	transport.TLSHandshakeTimeout = 1 * time.Second
	transport.ExpectContinueTimeout = 1 * time.Second
	transport.ResponseHeaderTimeout = 10 * time.Second
	transport.MaxIdleConns = 0
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return transport
}

func (s *downloaderService) Get(w io.Writer, ctx context.Context, config model.DownloadConfig) (int, *model.DownloadError, error) {
	if config.Compress {
		gzipWriter := gzip.NewWriter(w)

		defer func() {
			err := gzipWriter.Close()

			if err != nil {
				log.Print(err)
			}
		}()

		w = gzipWriter
	}

	if config.Url == "" {
		err := errors.New("url must be not empty")

		return 0, &model.DownloadError{
			ErrorState:   model.UrlNotValid,
			ErrorText:    err.Error(),
			ClientStatus: 0,
		}, err
	}

	client := new(http.Client)
	client.Timeout = config.Timeout

	if config.Proxy {
		s.configureProxy(client)
	} else {
		client.Transport = s.configureTransport(&http.Transport{})
	}

	req, err := http.NewRequestWithContext(ctx, "GET", config.Url, nil)

	if err != nil {

		return 0, &model.DownloadError{
			ErrorState:   model.InternalError,
			ErrorText:    err.Error(),
			ClientStatus: 0,
		}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, &model.DownloadError{
			ErrorState:   model.InternalHttpClientError,
			ErrorText:    err.Error(),
			ClientStatus: 0,
		}, err
	}

	headerWriter, isHeaderWriter := w.(http.ResponseWriter)

	if isHeaderWriter {
		headerWriter.WriteHeader(resp.StatusCode)
	}

	defer resp.Body.Close()

	if config.Sanitize.CleanMinimal {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return resp.StatusCode, &model.DownloadError{
				ErrorState:   model.InternalHttpClientError,
				ErrorText:    err.Error(),
				ClientStatus: resp.StatusCode,
			}, err
		}

		_, err = w.Write(gohtml.FormatBytes(s.sanitizer.SanitizeBytes(gohtml.FormatBytes(body))))

		if err != nil {
			return resp.StatusCode, &model.DownloadError{
				ErrorState:   model.SanitizerError,
				ErrorText:    err.Error(),
				ClientStatus: resp.StatusCode,
			}, err
		}
	}

	_, err = io.CopyN(w, resp.Body, 1024*1024*1024)

	if err != nil && err != io.EOF {
		return resp.StatusCode, &model.DownloadError{
			ErrorState:   model.InternalHttpClientError,
			ErrorText:    err.Error(),
			ClientStatus: resp.StatusCode,
		}, err
	}

	return resp.StatusCode, nil, nil
}

var DownloaderServiceInstance = new(downloaderService)

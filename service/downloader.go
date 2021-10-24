package service

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"download-delegator/model"
	ddError "download-delegator/model/errors"
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
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"
)

type downloaderService struct {
	ProxyFile  string
	proxyList  []model.ProxyItemConfig
	sanitizer  *bluemonday.Policy
	sanitizer2 *bluemonday.Policy
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

	s.sanitizer2 = bluemonday.NewPolicy()

	// We only allow <p> and <a href="">
	s.sanitizer2.AllowAttrs("name", "content", "property").OnElements("meta")
	s.sanitizer2.AllowElements("meta", "html", "head", "title")
	s.sanitizer2.SkipElementsContent("body")
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

func (s *downloaderService) configureProxy(client *http.Client, config model.DownloadConfig) {
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
	}, config)
}

func (s *downloaderService) configureTransport(transport *http.Transport, config model.DownloadConfig) *http.Transport {
	transport.TLSHandshakeTimeout = config.Timeout.TLSHandshakeTimeout
	transport.ExpectContinueTimeout = 1 * time.Second
	transport.MaxIdleConns = 0
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	transport.DialContext = (&net.Dialer{
		Timeout:   config.Timeout.DialTimeout,
		KeepAlive: -1,
	}).DialContext

	return transport
}

func (s *downloaderService) Get(w io.Writer, ctx context.Context, config model.DownloadConfig) (int, ddError.State, error) {
	select {
	case <-ctx.Done(): //context cancelled
		return 0, ddError.NoError, nil
	//case <-time.After(100 * time.Second): //timeout
	default:

	}

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

		return 0, ddError.UrlNotValid, err
	}

	client := new(http.Client)
	client.Timeout = config.Timeout.RequestTimeout
	defer func() {
		client.CloseIdleConnections()
	}()

	if config.Proxy {
		s.configureProxy(client, config)
	} else {
		client.Transport = s.configureTransport(&http.Transport{}, config)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", config.Url, nil)

	if err != nil {
		log.Print(err)

		return 0, ddError.InternalError, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, s.handleClientError(err), err
	}

	headerWriter, isHeaderWriter := w.(http.ResponseWriter)

	if isHeaderWriter {
		headerWriter.WriteHeader(resp.StatusCode)
	}

	defer resp.Body.Close()

	if config.Sanitize.CleanMinimal || config.Sanitize.CleanMinimal2 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Print(err)
			return resp.StatusCode, ddError.InternalHttpClientError, err
		}

		if config.Sanitize.CleanMinimal {
			_, err = w.Write(gohtml.FormatBytes(s.sanitizer.SanitizeBytes(gohtml.FormatBytes(body))))
		} else if config.Sanitize.CleanMinimal2 {
			_, err = w.Write(gohtml.FormatBytes(s.sanitizer2.SanitizeBytes(gohtml.FormatBytes(body))))
		}

		if err != nil {
			log.Print(err)
			return resp.StatusCode, ddError.SanitizerError, err
		}
	}

	_, err = io.CopyN(w, resp.Body, 1024*1024*1024)

	if err != nil && err != io.EOF {
		log.Print(err)
		return resp.StatusCode, ddError.InternalHttpClientError, err
	}

	if resp.StatusCode >= 400 {
		return resp.StatusCode, ddError.ClientNotSuccess, err
	}

	return resp.StatusCode, ddError.NoError, nil
}

func (s *downloaderService) handleClientError(err error) ddError.State {
	log.Print(err)

	err = unwrapErrorRecursive(err)

	if timeoutError, ok := err.(net.Error); ok && timeoutError.Timeout() {
		if strings.Contains(err.Error(), "dial tcp") {
			return ddError.DialTimeout
		} else if strings.Contains(err.Error(), "TLS handshake timeout") {
			return ddError.TlsTimeout
		} else {
			return ddError.Timeout
		}
	}

	if dnsError, ok := err.(*net.DNSError); ok && dnsError.Timeout() {
		return ddError.DnsTimeout
	}

	if dnsError, ok := err.(*net.DNSError); ok && !dnsError.Timeout() {
		return ddError.DnsNotResolved
	}

	if sysCallError, ok := err.(syscall.Errno); ok {
		if sysCallError == syscall.ECONNREFUSED {
			return ddError.ConnectionRefused
		}
		return ddError.SysCallGenericError
	}

	log.Print("client error: ", err)

	return ddError.InternalHttpClientError
}

func unwrapErrorRecursive(err error) error {
	newErr := errors.Unwrap(err)

	if newErr != nil {
		return unwrapErrorRecursive(newErr)
	}

	return err
}

var DownloaderServiceInstance = new(downloaderService)

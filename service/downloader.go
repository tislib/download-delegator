package service

import (
	"context"
	"crypto/tls"
	model2 "download-delegator/lib/parser/model"
	"download-delegator/model"
	ddError "download-delegator/model/errors"
	"download-delegator/service/transformers"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
)

type DownloaderService struct {
	proxyList    []model.ProxyItemConfig
	transformer  *TransformerService
	transformers []transformers.Transformer
}

func (s *DownloaderService) locateRandomProxy() *model.ProxyItemConfig {
	if len(s.proxyList) > 0 {
		randomIndex := rand.Intn(len(s.proxyList))
		return &s.proxyList[randomIndex]
	}

	return nil
}

func (s *DownloaderService) configureProxy(client *http.Client, config model.DownloadConfig) {
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

func (s *DownloaderService) configureTransport(transport *http.Transport, config model.DownloadConfig) *http.Transport {
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

func (s *DownloaderService) Get(ctx context.Context, config model.DownloadConfig) model.DownloadResponse {
	beginTime := time.Now()

	select {
	case <-ctx.Done(): //context cancelled
		return model.DownloadResponse{
			Url:        config.Url,
			Duration:   time.Now().Sub(beginTime),
			DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
			Error:      ddError.Timeout,
		}
	//case <-time.After(100 * time.Second): //timeout
	default:

	}

	if config.Url == "" {
		return model.DownloadResponse{
			Url:        config.Url,
			Duration:   time.Now().Sub(beginTime),
			DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
			Error:      ddError.UrlNotValid,
		}
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

		return model.DownloadResponse{
			Url:        config.Url,
			Duration:   time.Now().Sub(beginTime),
			DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
			Error:      ddError.InternalError,
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)

		return model.DownloadResponse{
			Url:        config.Url,
			Duration:   time.Now().Sub(beginTime),
			DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
			Error:      s.handleClientError(err),
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var err2 ddError.State

	if err != nil {
		log.Print(err)

		return model.DownloadResponse{
			Url:        config.Url,
			Duration:   time.Now().Sub(beginTime),
			DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
			Error:      ddError.InternalHttpClientError,
		}
	}

	if len(config.Transform) > 0 {
		body, err2 = s.transformer.Transform(body)

		if err2 != ddError.NoError {
			return model.DownloadResponse{
				Url:        config.Url,
				Duration:   time.Now().Sub(beginTime),
				DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
				Error:      err2,
			}
		}

		if err != nil {
			log.Print(err)

			return model.DownloadResponse{
				Url:        config.Url,
				Duration:   time.Now().Sub(beginTime),
				DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
				StatusCode: resp.StatusCode,
				Error:      ddError.WriterError,
			}
		}
	}

	errs := ddError.NoError

	if resp.StatusCode >= 400 {
		errs = ddError.ClientNotSuccess
	}

	return model.DownloadResponse{
		Url:        config.Url,
		Duration:   time.Now().Sub(beginTime),
		DurationMS: int(time.Now().Sub(beginTime) / time.Millisecond),
		StatusCode: resp.StatusCode,
		Content:    string(body),
		Error:      errs,
	}
}

func (s *DownloaderService) handleClientError(err error) ddError.State {
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

func (s *DownloaderService) InitTransformers(transformerConfigs []model2.TransformerConfig) {
	s.transformer = new(TransformerService)
	s.transformer.Init(transformerConfigs)
}

func unwrapErrorRecursive(err error) error {
	newErr := errors.Unwrap(err)

	if newErr != nil {
		return unwrapErrorRecursive(newErr)
	}

	return err
}

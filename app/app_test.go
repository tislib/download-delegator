package app

import (
	"context"
	"crypto/tls"
	client2 "download-delegator/client"
	"download-delegator/core/model"
	error2 "download-delegator/core/model/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func runServer() *App {
	app := new(App)

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

var app *App
var client *client2.DownloadDelegatorClient

func init() {
	app = runServer()
	client = new(client2.DownloadDelegatorClient)
	client.Init(client2.InitConfig{Addr: app.GetAddr()})

	log.SetLevel(log.TraceLevel)
}

func TestGetOnUrlWithoutProtocolAndGetProtocolSchemaError(t *testing.T) {
	log.Print(app.GetAddr())

	config := model.DownloadConfig{
		Url: "non-existing-domain-123.com",
	}

	response, err := client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.UnsupportedProtocolSchema, response.Error)
}

func TestGetOnUrlWithNonExistentDomainAndGetDnsError(t *testing.T) {
	log.Print(app.GetAddr())

	config := model.DownloadConfig{
		Url: "http://non-existing-domain-123.com",
	}

	response, err := client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.DnsNotResolved, response.Error)
}

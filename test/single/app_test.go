package single

import (
	"context"
	"download-delegator/core/model"
	error2 "download-delegator/core/model/errors"
	"download-delegator/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestGetOnUrlWithoutProtocolAndGetProtocolSchemaError(t *testing.T) {
	config := model.DownloadConfig{
		Url: "non-existing-domain-123.com",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.UnsupportedProtocolSchema, response.Error)
}

func TestGetOnUrlWithNonExistentDomainAndGetDnsError(t *testing.T) {
	config := model.DownloadConfig{
		Url: "http://non-existing-domain-123.com",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.DnsNotResolved, response.Error)
}

func TestGetLightServerCallingUrl(t *testing.T) {
	config := model.DownloadConfig{
		Url: "http://127.0.0.1:" + strconv.Itoa(test.LightTestServer.Port) + "/get-raw?data=test-data-123",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, error2.NoError, response.Error)
	assert.Equal(t, "test-data-123", response.Content)
}

func TestGetLightServerCallingUrlMediumPayload(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	defer func() {
		log.SetLevel(log.TraceLevel)
	}()

	config := model.DownloadConfig{
		Url: "http://127.0.0.1:" + strconv.Itoa(test.LightTestServer.Port) + "/get-raw?data=test-data-123&repeat=15000",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, error2.NoError, response.Error)
	assert.Equal(t, len("test-data-123")*15000, len(response.Content))
}

func TestGetLightServerCallingUrlLargePayload(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	defer func() {
		log.SetLevel(log.TraceLevel)
	}()

	config := model.DownloadConfig{
		Url: "http://127.0.0.1:" + strconv.Itoa(test.LightTestServer.Port) + "/get-raw?data=test-data-123&repeat=150000",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, error2.NoError, response.Error)
	assert.Equal(t, len("test-data-123")*150000, len(response.Content))
}

func TestGetLightServerTimeout(t *testing.T) {

	config := model.DownloadConfig{
		Url: "http://127.0.0.1:" + strconv.Itoa(test.LightTestServer.Port) + "/timeout?val=10000",
		Timeout: model.TimeoutConfig{
			TLSHandshakeTimeout: 100 * time.Millisecond,
			DialTimeout:         100 * time.Millisecond,
			RequestTimeout:      200 * time.Millisecond,
		},
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.Timeout, response.Error)
	assert.Equal(t, "", response.Content)
}

func TestConnectionRefused(t *testing.T) {
	config := model.DownloadConfig{
		Url: "http://127.0.0.1:12345",
	}

	response, err := test.Client.Get(context.TODO(), config)

	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, error2.ConnectionRefused, response.Error)
	assert.Equal(t, "", response.Content)
}

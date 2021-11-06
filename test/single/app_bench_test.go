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
)

func BenchmarkSimpleGet(b *testing.B) {
	log.SetLevel(log.InfoLevel)
	defer func() {
		log.SetLevel(log.TraceLevel)
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			config := model.DownloadConfig{
				Url: "http://127.0.0.1:" + strconv.Itoa(test.LightTestServer.Port) + "/get-raw?data=test-data-123&repeat=15000",
			}

			response, err := test.Client.Get(context.TODO(), config)

			assert.Nil(b, err)
			assert.Equal(b, 200, response.StatusCode)
			assert.Equal(b, error2.NoError, response.Error)
			assert.Equal(b, len("test-data-123")*15000, len(response.Content))
		}
	})
}

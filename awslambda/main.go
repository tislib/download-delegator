package main

import (
	"context"
	"download-delegator/app"
	"download-delegator/model"
	"download-delegator/model/errors"
	"download-delegator/service"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) ([]model.DownloadResponse, error) {
	timeCalc := new(app.TimeCalc)
	timeCalc.Init("bulk-download")

	wg := new(sync.WaitGroup)

	var config model.BulkDownloadConfig

	err := json.Unmarshal([]byte(request.Body), &config)

	if err != nil {
		log.Panic(err)
	}

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	log.Print("accepted config: ", config)

	var maxConcurrencySemaphore = make(chan int, config.MaxConcurrency)

	var result []model.DownloadResponse

	var counter int32
	for indexX, itemX := range config.Url {
		wg.Add(1)
		maxConcurrencySemaphore <- 1

		go func(index int, item string) {
			atomic.AddInt32(&counter, 1)

			defer func() {
				atomic.AddInt32(&counter, -1)
				log.Print("concurrency level: ", strconv.Itoa(int(counter)))

				if err := recover(); err != nil {
					log.Print("recovered error: ", err, string(debug.Stack()))
				}
				wg.Done()
				<-maxConcurrencySemaphore
			}()

			select {
			case <-ctx.Done():
				return
			default:

			}

			log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"begin bulk download index/url: ", index, item)

			var resItem model.DownloadResponse

			var downloaderServiceInstance = new(service.DownloaderService)
			downloaderServiceInstance.ConfigureTransformers(config.Transform)
			downloaderServiceInstance.EnableProxy(config.Proxy)
			downloaderServiceInstance.ConfigureTimeout(config.Timeout)

			for i := 0; i < config.RetryCount; i++ {
				resItem = downloaderServiceInstance.Get(ctx, item)

				if resItem.Error != errors.NoError {
					log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"end bulk download index/url: ", index, item, resItem.StatusCode, len(resItem.Content), resItem.DurationMS)
					break
				}
			}

			timeCalc.Step()
			log.Print("sending to chan", counter)
			result = append(result, resItem)
			log.Print("sent to chan", counter)
		}(indexX, itemX)
	}

	wg.Wait()

	return result, nil
}

func main() {
	lambda.Start(HandleRequest)
}

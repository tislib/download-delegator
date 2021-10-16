package main

import (
	"bytes"
	"context"
	"download-delegator/app"
	"download-delegator/model"
	"download-delegator/service"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
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

	result = append(result, model.DownloadResponse{
		Url:           "dummy-url",
		Content:       "",
		StatusCode:    0,
		Duration:      0,
		DurationMS:    0,
		DownloadError: nil,
		Index:         0,
		Retried:       0,
	})

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

			downloadConfig := model.DownloadConfig{
				Url:      item,
				Compress: false,
				Sanitize: config.Sanitize,
				Proxy:    config.Proxy,
				Timeout:  config.Timeout,
			}

			var buf = bytes.Buffer{}

			beginTime := time.Now()

			var resItem model.DownloadResponse

			for i := 0; i < config.RetryCount; i++ {
				statusCode, downloadErr, err := service.DownloaderServiceInstance.Get(&buf, ctx, downloadConfig)

				duration := time.Now().Sub(beginTime)

				localResItem := model.DownloadResponse{
					Url:           item,
					Index:         index,
					StatusCode:    statusCode,
					Content:       buf.String(),
					DownloadError: downloadErr,
					Duration:      duration,
					Retried:       i,
					DurationMS:    int(duration / time.Millisecond),
				}

				resItem = localResItem

				if err != nil {
					log.Print("["+(strconv.Itoa(len(config.Url)))+"/"+strconv.Itoa(index)+"]"+"end bulk download index/url: ", index, item, statusCode, len(resItem.Content), int(duration/time.Millisecond))
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
	service.DownloaderServiceInstance.ConfigureSanitizer()

	lambda.Start(HandleRequest)
}

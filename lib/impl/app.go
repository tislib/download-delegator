package service

import (
	"context"
	"download-delegator/core/model"
)

type AppServiceImpl struct {
}

func (s AppServiceImpl) Get(ctx context.Context, config model.DownloadConfig) (*model.DownloadResponse, error) {
	downloaderService := new(DownloaderService)
	downloaderService.ConfigureTransformers(config.Transform)
	downloaderService.timeout = config.Timeout

	downloadResponse := downloaderService.Get(ctx, config.Url)

	return &downloadResponse, nil
}

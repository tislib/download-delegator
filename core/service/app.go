package service

import (
	"context"
	"download-delegator/core/model"
)

type AppService interface {
	Get(ctx context.Context, config model.DownloadConfig) (*model.DownloadResponse, error)
}

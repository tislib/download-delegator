package model

import (
	error2 "download-delegator/core/model/errors"
	"download-delegator/lib/parser/model"
	"time"
)

type OutputForm string

const (
	JsonOutput OutputForm = "json"
)

type BulkDownloadConfig struct {
	Url            []string
	Compression    Compression
	Proxy          bool
	Timeout        TimeoutConfig
	Transform      []model.TransformerConfig
	OutputForm     OutputForm
	MaxConcurrency int
	RetryCount     int
}

type DownloadResponse struct {
	Url        string
	Content    string
	StatusCode int
	Duration   time.Duration
	DurationMS int
	Error      error2.State
	Index      int
	Retried    int
}

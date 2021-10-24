package model

import (
	error2 "download-delegator/model/errors"
	"time"
)

type BulkWhoisConfig struct {
	Domains        []string
	Compression    Compression
	Timeout        time.Duration
	OutputForm     OutputForm
	MaxConcurrency int
	RetryCount     int
}

type WhoisResponse struct {
	Domain     string
	Response   string
	Duration   time.Duration
	DurationMS int
	Error      error2.State
	Index      int
	Retried    int
}

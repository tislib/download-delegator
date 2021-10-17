package model

import (
	"time"
)

type BulkWhoisConfig struct {
	Domains        []string
	Compress       bool
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
	Error      *Error
	Index      int
	Retried    int
}

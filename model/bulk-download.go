package model

import "time"

type OutputForm string

const (
	JsonOutput OutputForm = "json"
)

type BulkDownloadConfig struct {
	Url            []string
	Compression    Compression
	Proxy          bool
	Timeout        TimeoutConfig
	Sanitize       SanitizeConfig
	OutputForm     OutputForm
	MaxConcurrency int
	RetryCount     int
}

type DownloadResponse struct {
	Url           string
	Content       string
	StatusCode    int
	Duration      time.Duration
	DurationMS    int
	DownloadError *Error
	Index         int
	Retried       int
}

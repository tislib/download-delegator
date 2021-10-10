package model

import "time"

type OutputForm string

const (
	JsonOutput OutputForm = "json"
)

type BulkDownloadConfig struct {
	Url        []string
	Compress   bool
	Proxy      bool
	Timeout    time.Duration
	Sanitize   SanitizeConfig
	OutputForm OutputForm
}

type DownloadResponse struct {
	Url           string
	Content       string
	StatusCode    int
	Duration      time.Duration
	DownloadError *DownloadError
	Index         int
}

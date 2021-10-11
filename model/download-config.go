package model

import "time"

type DownloadConfig struct {
	Url      string
	Compress bool
	Proxy    bool
	Timeout  TimeoutConfig
	Sanitize SanitizeConfig
}

type TimeoutConfig struct {
	TLSHandshakeTimeout time.Duration
	DialTimeout         time.Duration
	RequestTimeout      time.Duration
}

type SanitizeConfig struct {
	CleanMinimal  bool
	CleanMinimal2 bool
}

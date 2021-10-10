package model

import "time"

type DownloadConfig struct {
	Url      string
	Compress bool
	Proxy    bool
	Timeout  time.Duration
	Sanitize SanitizeConfig
}

type SanitizeConfig struct {
	CleanMinimal  bool
	CleanMinimal2 bool
}

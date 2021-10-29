package model

import (
	"download-delegator/lib/parser/model"
	"time"
)

type DownloadConfig struct {
	Url       string
	Proxy     bool
	Timeout   TimeoutConfig
	Transform []model.TransformerConfig
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

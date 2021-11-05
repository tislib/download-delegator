package service

import (
	"download-delegator/model"
	error2 "download-delegator/model/errors"
	whois "github.com/brimstone/golang-whois"
	log "github.com/sirupsen/logrus"
	"time"
)

type whoisService struct {
}

func (s *whoisService) Get(domain string, timeout time.Duration) model.WhoisResponse {
	beginTime := time.Now()

	if domain == "" {
		return model.WhoisResponse{
			Error: error2.DomainNotValid,
		}
	}

	data, err := whois.GetWhoisTimeout(domain, timeout)

	if err != nil {
		log.Print(err)
	}

	duration := time.Now().Sub(beginTime)

	return model.WhoisResponse{
		Domain:     domain,
		Response:   data.String(),
		Duration:   duration,
		DurationMS: int(duration / time.Millisecond),
		Error:      error2.WhoisError,
	}
}

var WhoisServiceInstance = new(whoisService)

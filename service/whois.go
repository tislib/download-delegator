package service

import (
	"download-delegator/model"
	whois "github.com/brimstone/golang-whois"
	"log"
	"time"
)

type whoisService struct {
}

func (s *whoisService) Get(domain string, timeout time.Duration) model.WhoisResponse {
	beginTime := time.Now()

	if domain == "" {
		return model.WhoisResponse{
			Error: model.DomainNotValid.Error(),
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
		Error:      model.WhoisError.ErrorWithError(err),
	}
}

var WhoisServiceInstance = new(whoisService)

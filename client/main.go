package client

import (
	"bytes"
	"context"
	"download-delegator/core/model"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type DownloadDelegatorClient struct {
	addr string
}

func (receiver *DownloadDelegatorClient) Init(initConfig InitConfig) {
	receiver.addr = initConfig.Addr
}

func (receiver DownloadDelegatorClient) Get(ctx context.Context, config model.DownloadConfig) (*model.DownloadResponse, error) {
	url := "http://" + receiver.addr + "/get"

	requestData, err := json.Marshal(config)

	resp, err := http.Post(url, "application/json", bytes.NewReader(requestData))

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)

	log.Trace("GET Response: ", string(responseData))

	var downloadResponse = new(model.DownloadResponse)

	err = json.Unmarshal(responseData, downloadResponse)

	if err != nil {
		return downloadResponse, err
	}

	if resp.StatusCode != 200 {
		return downloadResponse, errors.New("operational error: " + resp.Status)
	}

	return downloadResponse, err
}

package main

import (
	"bytes"
	"crypto/tls"
	"download-delegator/model"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const url = "https://localhost:8000"

var httpVersion = flag.Int("version", 2, "HTTP version")

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	bulkDownload := new(model.BulkDownloadConfig)

	bulkDownload.MaxConcurrency = 500
	bulkDownload.Sanitize = model.SanitizeConfig{CleanMinimal2: true}
	bulkDownload.OutputForm = model.JsonOutput
	bulkDownload.RetryCount = 10
	bulkDownload.Timeout = model.TimeoutConfig{
		TLSHandshakeTimeout: time.Second * 10,
		DialTimeout:         time.Second * 10,
		RequestTimeout:      time.Second * 10,
	}

	N := 3

	data, _ := os.ReadFile("/Users/taleh/Downloads/domains/domains100k.sample.text")

	dataStr := string(data)

	var i = N
	for _, domain := range strings.Split(dataStr, "\n") {
		i--
		bulkDownload.Url = append(bulkDownload.Url, "http://"+domain)

		if i == 0 {
			break
		}
	}

	beginTime := time.Now()

	m, b := bulkDownload, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	r, e := http.NewRequest("POST", "https://127.0.0.1:8234/bulk", b)
	if e != nil {
		panic(e)
	}

	resp, e := new(http.Client).Do(r)

	duration := time.Now().Sub(beginTime)

	if e != nil {
		log.Panic(e)
	}

	var response []model.DownloadResponse

	err := json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Panic(err)
	}

	var downloadErrorStats = make(map[model.DownloadErrorState]int)

	for _, item := range response {
		if item.DownloadError != nil {
			//log.Print(item.DownloadError.ErrorText)
			downloadErrorStats[item.DownloadError.ErrorState]++
		} else {
			downloadErrorStats["ok"]++
		}
	}

	log.Print(response)

	log.Println("duration: ", duration)
	log.Println("rps: ", int(time.Second)/(int(duration)/N))
	log.Print(downloadErrorStats)
	log.Print("T: ", bulkDownload.Timeout.RequestTimeout)
	log.Print("N: ", N)
	log.Print("C: ", bulkDownload.MaxConcurrency)
	log.Print("R: ", bulkDownload.RetryCount)
}

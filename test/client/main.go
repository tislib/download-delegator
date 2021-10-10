package main

import (
	"bytes"
	"crypto/tls"
	"download-delegator/model"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
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
	bulkDownload.Timeout = time.Second * 100

	N := 1000

	for i := 0; i < N; i++ {
		bulkDownload.Url = append(bulkDownload.Url, "https://static.tisserv.net/")
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

	data, err := ioutil.ReadAll(resp.Body)

	log.Println(string(data), err)

	log.Println("duration: ", duration)
	log.Println("rps: ", int(duration/time.Millisecond)/N)
}

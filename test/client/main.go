package main

import (
	"bytes"
	"crypto/tls"
	model2 "download-delegator/lib/parser/model"
	"download-delegator/model"
	error2 "download-delegator/model/errors"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

import _ "net/http/pprof"

const url = "https://localhost:8000"

var httpVersion = flag.Int("version", 2, "HTTP version")

func main() {
	log.Print(int64(time.Second))

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	bulkDownload := new(model.BulkDownloadConfig)

	bulkDownload.MaxConcurrency = 150
	bulkDownload.OutputForm = model.JsonOutput
	bulkDownload.RetryCount = 1
	bulkDownload.Timeout = model.TimeoutConfig{
		TLSHandshakeTimeout: time.Second * 3,
		DialTimeout:         time.Second * 3,
		RequestTimeout:      time.Second * 10,
	}

	bulkDownload.Transform = append(bulkDownload.Transform, model2.TransformerConfig{
		Type: model2.Sanitize,
	})

	bulkDownload.Transform = append(bulkDownload.Transform, model2.TransformerConfig{
		Type: model2.HtmlFormat,
	})

	//bulkDownload.Transform = append(bulkDownload.Transform, model.TransformerConfig{
	//	Type: model.ScriptTengo,
	//	Parameters: `output := input`,
	//})

	bulkDownload.Transform = append(bulkDownload.Transform, model2.TransformerConfig{
		Type: model2.ScriptJs,
		Parameters: `


res = lib.ParseHtml(input)

output = JSON.stringify({
    "res": res.FindSingle("title").Text()
})

		`,
	})

	N := 6000

	data, _ := os.ReadFile("/Users/taleh/Downloads/domains/domains100k.sample.text")

	dataStr := string(data)

	var i = N
	for _, domain := range strings.Split(dataStr, "\n") {
		i--
		bulkDownload.Url = append(bulkDownload.Url, "https://www.alexa.com/siteinfo/"+domain)

		if i == 0 {
			break
		}
	}

	beginTime := time.Now()

	m, b := bulkDownload, new(bytes.Buffer)

	json.NewEncoder(b).Encode(m)

	//payload := string(b.Bytes())
	//log.Print(payload)

	r, e := http.NewRequest("POST", "http://localhost:8234/bulk", b)
	if e != nil {
		panic(e)
	}

	dClient := new(http.Client)
	dClient.Timeout = 1 * time.Hour

	resp, e := dClient.Do(r)

	duration := time.Now().Sub(beginTime)

	if e != nil {
		log.Panic(e)
	}

	var response []model.DownloadResponse

	err := json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Panic(err)
	}

	var downloadErrorStats = make(map[error2.State]int)

	for _, item := range response {
		if item.Error != error2.NoError {
			//log.Print(item.DownloadError.ErrorText)
			downloadErrorStats[item.Error]++
		} else {
			downloadErrorStats["ok"]++
		}
	}

	//log.Print(response)

	log.Print(response[0].Content)
	log.Println("duration: ", duration)
	log.Println("rps: ", int(time.Second)/(int(duration)/N))
	log.Print(downloadErrorStats)
	log.Print("T: ", bulkDownload.Timeout.RequestTimeout)
	log.Print("N: ", N)
	log.Print("C: ", bulkDownload.MaxConcurrency)
	log.Print("R: ", bulkDownload.RetryCount)
}

//func main() {
//
//	N := 10000
//
//	data, _ := os.ReadFile("/Users/taleh/Downloads/domains/com100k.domains.txt")
//
//	dataStr := string(data)
//
//	var i = N
//	sem := semaphore.NewWeighted(100)
//	wg := sync.WaitGroup{}
//
//	timeCalc := new(app.TimeCalc)
//	timeCalc.Init("time")
//
//	for _, domain := range strings.Split(dataStr, "\n") {
//		err := sem.Acquire(context.TODO(), 1)
//		wg.Add(1)
//
//		if err != nil {
//			log.Panic(err)
//		}
//
//		i--
//		go func() {
//			//r := &net.Resolver{
//			//	PreferGo: true,
//			//	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
//			//		d := net.Dialer{
//			//			Timeout: time.Millisecond * time.Duration(10000),
//			//		}
//			//		return d.DialContext(ctx, network, "8.8.8.8:53")
//			//	},
//			//}
//			//ip, err := r.LookupHost(context.Background(), domain)
//
//			//log.Println(domain)
//			data, err := whois.GetWhois(domain)
//
//			if err != nil {
//				log.Print(err)
//			}
//
//			if i%1000 == 0 {
//				log.Println(data)
//			}
//			sem.Release(1)
//			wg.Done()
//			timeCalc.Step()
//		}()
//
//		if i == 0 {
//			break
//		}
//	}
//
//	wg.Wait()
//
//}

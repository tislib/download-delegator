package main

import (
	"context"
	"crypto/tls"
	"download-delegator/app"
	"golang.org/x/net/http2"
	"golang.org/x/sync/semaphore"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	timeCalcR := new(app.TimeCalc)
	timeCalcR.Init("Request")
	timeCalcO := new(app.TimeCalc)
	timeCalcO.Init("OK")

	timeCalcE := new(app.TimeCalc)
	timeCalcE.Init("ERROR")

	client := new(http.Client)
	client.Timeout = time.Second * 100

	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var sem = semaphore.NewWeighted(int64(1000))

	for {
		sem.Acquire(context.TODO(), 1)
		go func() {
			timeCalcR.Step()
			runC(client, timeCalcO, timeCalcE)
			sem.Release(1)
		}()

		time.Sleep(1 * time.Millisecond)
	}

}

func runC(client *http.Client, timeCalcO *app.TimeCalc, timeCalcE *app.TimeCalc) {
	resp, err := client.Get("https://tisserv.net:8234/get?url=https://tap.az/elanlar/elektronika/audio-video/20571375")

	if err != nil {
		//log.Print(err)
		timeCalcE.Step()
		return
	}

	defer resp.Body.Close()

	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)
	log.Print(len(buf.String()))

	if resp.StatusCode == 200 {
		timeCalcO.Step()
	} else {
		timeCalcE.Step()
	}
}

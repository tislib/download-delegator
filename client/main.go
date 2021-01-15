package main

import (
	"crypto/tls"
	"download-delegator/app"
	"golang.org/x/net/http2"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	timeCalcO := new(app.TimeCalc)
	timeCalcO.Init("OK")

	timeCalcE := new(app.TimeCalc)
	timeCalcE.Init("ERROR")

	client := new(http.Client)

	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for {
		go func() {
			runC(client, timeCalcO, timeCalcE)
		}()

		time.Sleep(10 * time.Millisecond)
	}

}

func runC(client *http.Client, timeCalcO *app.TimeCalc, timeCalcE *app.TimeCalc) {
	resp, err := client.Get("https://tisserv.net:8234/get?url=https://www.gittigidiyor.com/saglik-medikal/medikal-sarf-malzemeleri")

	if err != nil {
		//log.Print(err)
		timeCalcE.Step()
		return
	}

	defer resp.Body.Close()

	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)
	//log.Print(n)

	if resp.StatusCode == 200 {
		timeCalcO.Step()
	} else {
		timeCalcE.Step()
	}
}

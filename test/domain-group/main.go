package main

import (
	"bufio"
	"download-delegator/app"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("/Users/taleh/Downloads/domains/domains.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var tldCounter = make(map[string]int64)

	timeCalc := new(app.TimeCalc)
	timeCalc.Init("dd")

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ".")
		if len(parts) != 2 {
			continue
		}
		tld := parts[1]

		tldCounter[tld]++
		timeCalc.Step()
	}

	log.Println(tldCounter)

}

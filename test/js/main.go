package main

import (
	"download-delegator/service/transformers"
	"log"
	"os"
)

func main() {
	var input = []byte("<p>hello world 123</p>")

	transformer := new(transformers.JsTransformer)

	script, _ := os.ReadFile("/Users/taleh/Projects/UniersalDataPlatform/download-delegator/test/js/test.js")

	transformer.Init(string(script))

	out, e2 := transformer.Transform(input)

	log.Println(string(out), e2)
}

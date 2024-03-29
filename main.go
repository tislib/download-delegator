package main

import (
	appPackage "download-delegator/app"
	"download-delegator/core/model"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	configFile, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Panicln(err)
	}

	var config model.Config

	_, err = toml.Decode(string(configFile), &config)

	if err != nil {
		log.Panicln(err)
	}

	app := new(appPackage.App)

	app.Init(config)

	app.Version = "1.0.3"

	app.Run()
}

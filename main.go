package main

import (
	appPackage "download-delegator/app"
	"os"
)

func main() {
	app := new(appPackage.App)

	app.Addr = ":8234"

	app.CertFile = os.Args[1]
	app.KeyFile = os.Args[2]
	app.ProxyFile = os.Args[3]

	app.Run()
}

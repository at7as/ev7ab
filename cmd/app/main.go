package main

import (
	"flag"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lib"
)

func main() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&lib.ExampleSimple{}, *cfgFile)

}

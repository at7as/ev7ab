package main

import (
	"flag"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lib"
)

func main() {

	labFile := flag.String("lab", "./ev.lab", "path to lab file")

	flag.Parse()

	app.Run(&lib.Example1{}, *labFile)

}

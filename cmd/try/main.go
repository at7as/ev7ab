package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/at7as/ev7ab/pkg/lib"
)

func main() {

	f, err := os.Open("./test/example_simple/ev.lab")
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	d, err := zlib.NewReader(f)
	if err != nil {
		log.Panicln(err)
	}
	defer d.Close()

	b, err := io.ReadAll(d)
	if err != nil {
		log.Panicln(err)
	}

	l := lab.New(&lib.ExampleSimple{})

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	fmt.Println(l.Value([]float64{0.1, 0.2}))
	fmt.Println(l.Volume([]float64{0.1, 0.2}))

}

// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lib

import (
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/fogleman/gg"
)

type ExampleDigits struct {
	ls *labelSet
	is *imageSet
}

func (p *ExampleDigits) Load(setup map[string]string) error {

	ls, err := newLabelSet("./test/example_digits/t10k-labels-idx1-ubyte.gz")
	if err != nil {
		log.Panicln(err)
	}
	p.ls = ls
	p.ls.items = p.ls.items[:1000]

	is, err := newImageSet("./test/example_digits/t10k-images-idx3-ubyte.gz")
	if err != nil {
		log.Panicln(err)
	}
	p.is = is
	p.is.images = p.is.images[:1000]

	return nil
}

func (p *ExampleDigits) Setup(key, value string) error {

	return nil
}

func (p *ExampleDigits) Produce(n lab.Next, _ lab.Next, _ []float64) []float64 {

	out := 0
	for i, v := range p.is.images {
		if int(math.Round(n(v)[0]*10.0)) == p.ls.items[i] {
			out++
		}
	}

	return []float64{float64(out) / float64(len(p.is.images))}
}

func (p *ExampleDigits) Validate(r []float64) bool {

	if r[0] < 0.1 {
		return false
	}

	return true
}

func (p *ExampleDigits) Compare(a, b []float64) bool {

	return a[0] > b[0]
}

func (p *ExampleDigits) Best(v []float64) string {

	best := "---"
	if len(v) > 0 {
		best = fmt.Sprintf("%.1f%%", v[0]*100.0)
	}

	return best
}

func (p *ExampleDigits) Goal(v []float64) bool {

	if v[0] > 0.95 {
		return true
	}

	return false
}

func ExampleDigitsApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleDigits{}, *cfgFile, false)

}

func ExampleDigitsTestDraw() {

	ls, err := newLabelSet("./test/example_digits/t10k-labels-idx1-ubyte.gz")
	if err != nil {
		log.Panicln(err)
	}

	is, err := newImageSet("./test/example_digits/t10k-images-idx3-ubyte.gz")
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(ls.items[:10])

	image := gg.NewContext(280, 28)
	for i := range 10 {
		digit := is.images[i]
		for y := range 28 {
			for x := range 28 {
				image.SetRGB(digit[x+(y*28)], digit[x+(y*28)], digit[x+(y*28)])
				image.SetPixel(x+i*28, y)
			}
		}
	}
	image.SavePNG("./test/example_digits/test_draw.png")

}

func newImageSet(path string) (*imageSet, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	s := &imageSet{header: imageSetHeader{}}
	if err := binary.Read(gz, binary.BigEndian, &s.header); err != nil {
		return nil, err
	}
	s.raw = make([]byte, s.header.Num*s.header.Rows*s.header.Columns)
	if err := binary.Read(gz, binary.BigEndian, &s.raw); err != nil {
		return nil, err
	}
	if s.header.Magic != 0x00000803 {
		return nil, fmt.Errorf("incorrect format")
	}
	s.parse()

	return s, nil
}

func newLabelSet(path string) (*labelSet, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	s := &labelSet{header: labelSetHeader{}}
	if err := binary.Read(gz, binary.BigEndian, &s.header); err != nil {
		return nil, err
	}
	s.raw = make([]byte, s.header.Num)
	if err := binary.Read(gz, binary.BigEndian, &s.raw); err != nil {
		return nil, err
	}
	if s.header.Magic != 0x00000801 {
		return nil, fmt.Errorf("incorrect format")
	}
	s.parse()

	return s, nil
}

type imageSet struct {
	header imageSetHeader
	raw    []byte
	images [][]float64
}

type imageSetHeader struct {
	Magic   int32
	Num     int32
	Rows    int32
	Columns int32
}

type labelSet struct {
	header labelSetHeader
	raw    []byte
	items  []int
}

type labelSetHeader struct {
	Magic int32
	Num   int32
}

func (s *imageSet) parse() {

	s.images = make([][]float64, s.header.Num)
	for i := range s.header.Num {
		s.images[i] = make([]float64, s.header.Rows*s.header.Columns)
		for ii, v := range s.raw[i*s.header.Rows*s.header.Columns : (i+1)*s.header.Rows*s.header.Columns] {
			s.images[i][ii] = float64(int(v)) / 255.0
		}
	}

}

func (s *labelSet) parse() {

	s.items = make([]int, s.header.Num)
	for i := range s.header.Num {
		s.items[i] = int(s.raw[i])
	}

}

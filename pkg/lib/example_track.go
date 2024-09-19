package lib

import (
	"bufio"
	"compress/zlib"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/fogleman/gg"
)

type ExampleTrack struct {
	d TrackData
	t [][]float64
}

func (p *ExampleTrack) Load(setup map[string]string) error {

	f, err := os.Open("./test/example_track/track.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var td TrackData
	if err = json.Unmarshal(b, &td); err != nil {
		return err
	}

	t := [][]float64{}

	l := len(td.Data) / 2

	for n := range l {

		tt := []float64{
			distance(td.Data[index(n*2, len(td.Data))][0], td.Data[index(n*2, len(td.Data))][1], td.Data[index(n*2+2, len(td.Data))][0], td.Data[index(n*2+2, len(td.Data))][1]),
			distance(td.Data[index(n*2+1, len(td.Data))][0], td.Data[index(n*2+1, len(td.Data))][1], td.Data[index(n*2+3, len(td.Data))][0], td.Data[index(n*2+3, len(td.Data))][1]),
			distance(td.Data[index(n*2+2, len(td.Data))][0], td.Data[index(n*2+2, len(td.Data))][1], td.Data[index(n*2+4, len(td.Data))][0], td.Data[index(n*2+4, len(td.Data))][1]),
			distance(td.Data[index(n*2+3, len(td.Data))][0], td.Data[index(n*2+3, len(td.Data))][1], td.Data[index(n*2+5, len(td.Data))][0], td.Data[index(n*2+5, len(td.Data))][1]),
			distance(td.Data[index(n*2+4, len(td.Data))][0], td.Data[index(n*2+4, len(td.Data))][1], td.Data[index(n*2+6, len(td.Data))][0], td.Data[index(n*2+6, len(td.Data))][1]),
			distance(td.Data[index(n*2+5, len(td.Data))][0], td.Data[index(n*2+5, len(td.Data))][1], td.Data[index(n*2+7, len(td.Data))][0], td.Data[index(n*2+7, len(td.Data))][1]),
			distance(td.Data[index(n*2+6, len(td.Data))][0], td.Data[index(n*2+6, len(td.Data))][1], td.Data[index(n*2+8, len(td.Data))][0], td.Data[index(n*2+8, len(td.Data))][1]),
			distance(td.Data[index(n*2+7, len(td.Data))][0], td.Data[index(n*2+7, len(td.Data))][1], td.Data[index(n*2+9, len(td.Data))][0], td.Data[index(n*2+9, len(td.Data))][1]),
		}
		tm := slices.Max(tt)
		for i, v := range tt {
			tt[i] = v / tm
		}

		t = append(t, tt)

	}

	p.d = td
	p.t = t

	return nil
}

func (p *ExampleTrack) Setup(key, value string) error {

	return nil
}

func (p *ExampleTrack) Produce(n lab.Next, op lab.Next) []float64 {

	tt := 0.5
	x0 := 0.0
	y0 := 0.0
	x1 := 0.0
	y1 := 0.0
	for i, in := range p.t {
		tt = n(append(in, tt))[0]
		if i == len(p.t)-2 {
			x0, y0 = lerp(p.d.Data[len(p.d.Data)-4][0], p.d.Data[len(p.d.Data)-4][1], p.d.Data[len(p.d.Data)-3][0], p.d.Data[len(p.d.Data)-3][1], tt)
		}
		if i == len(p.t)-1 {
			x1, y1 = lerp(p.d.Data[len(p.d.Data)-2][0], p.d.Data[len(p.d.Data)-2][1], p.d.Data[len(p.d.Data)-1][0], p.d.Data[len(p.d.Data)-1][1], tt)
		}
	}
	x2 := 0.0
	y2 := 0.0
	out := 180.0
	for i, in := range p.t {
		tt = n(append(in, tt))[0]
		x2, y2 = lerp(p.d.Data[i*2][0], p.d.Data[i*2][1], p.d.Data[i*2+1][0], p.d.Data[i*2+1][1], tt)
		out = min(out, deg(angle(x0, y0, x1, y1, x2, y2)))
		x0 = x1
		y0 = y1
		x1 = x2
		y1 = y2
	}

	return []float64{out}
}

func (p *ExampleTrack) Compare(a, b []float64) bool {

	return a[0] > b[0]
}

func (p *ExampleTrack) Validate(r []float64) bool {

	if r[0] < 90.0 {
		return false
	}

	return true
}

func (p *ExampleTrack) Best(v []float64) string {

	best := "---"
	if len(v) > 0 {
		best = fmt.Sprintf("%.1f°", v[0])
	}

	return best
}

func (p *ExampleTrack) Goal(v []float64) bool {

	return false
}

func ExampleTrackApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleTrack{}, *cfgFile, true)

}

func ExampleTrackTry(id int) {

	f, err := os.Open("./test/example_track/ev.lab")
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

	l := lab.New(&ExampleTrack{}, false)

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	track := &ExampleTrack{}
	if err := track.Load(nil); err != nil {
		log.Panicln(err)
	}

	ll := [][2]float64{}
	cl := [][2]float64{}
	rl := [][2]float64{}

	for n := range len(track.d.Data) / 2 {
		ll = append(ll, [2]float64{track.d.Data[n*2][0], track.d.Data[n*2][1]})
		rl = append(rl, [2]float64{track.d.Data[n*2+1][0], track.d.Data[n*2+1][1]})
	}

	tt := 0.5
	for _, in := range track.t {
		tt = l.ProjectValue(id, append(in, tt))[0]
	}
	for i, in := range track.t {
		tt = l.ProjectValue(id, append(in, tt))[0]
		x, y := lerp(ll[i][0], ll[i][1], rl[i][0], rl[i][1], tt)
		cl = append(cl, [2]float64{x, y})
	}

	size := 1024
	mx := 0.0
	my := 0.0
	for i := range len(cl) {
		mx = min(mx, ll[i][0])
		my = min(my, ll[i][1])
		mx = min(mx, rl[i][0])
		my = min(my, rl[i][1])
		mx = min(mx, cl[i][0])
		my = min(my, cl[i][1])
	}
	mx = math.Abs(mx)
	my = math.Abs(my)
	for i := range len(cl) {
		ll[i][0] += mx
		ll[i][1] += my
		ll[i][0] *= 40
		ll[i][1] *= 40
		ll[i][0] += 50
		ll[i][1] += 50
		rl[i][0] += mx
		rl[i][1] += my
		rl[i][0] *= 40
		rl[i][1] *= 40
		rl[i][0] += 50
		rl[i][1] += 50
		cl[i][0] += mx
		cl[i][1] += my
		cl[i][0] *= 40
		cl[i][1] *= 40
		cl[i][0] += 50
		cl[i][1] += 50
	}

	image := gg.NewContext(size, size)
	image.DrawRectangle(0, 0, float64(size), float64(size))
	image.SetRGB(1, 1, 1)
	image.Fill()
	for i, p := range ll {
		if i != len(ll)-1 {
			image.DrawLine(p[0], p[1], ll[i+1][0], ll[i+1][1])
		} else {
			image.DrawLine(p[0], p[1], ll[0][0], ll[0][1])
		}
	}
	image.SetRGB(1, 0, 0)
	image.Stroke()
	for i, p := range rl {
		if i != len(rl)-1 {
			image.DrawLine(p[0], p[1], rl[i+1][0], rl[i+1][1])
		} else {
			image.DrawLine(p[0], p[1], rl[0][0], rl[0][1])
		}
	}
	image.SetRGB(0, 0, 1)
	image.Stroke()
	for i, p := range cl {
		if i != len(cl)-1 {
			image.DrawLine(p[0], p[1], cl[i+1][0], cl[i+1][1])
		} else {
			image.DrawLine(p[0], p[1], cl[0][0], cl[0][1])
		}
	}
	image.SetRGB(0, 1, 0)
	image.Stroke()
	image.SavePNG("./test/example_track/track_result.png")

}

func ExampleTrackDraw() error {

	f, err := os.Open("./test/example_track/track.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var td TrackData
	if err = json.Unmarshal(b, &td); err != nil {
		return err
	}

	size := 1024
	mx := 0.0
	my := 0.0
	for _, p := range td.Data {
		mx = min(mx, p[0])
		my = min(my, p[1])
	}
	mx = math.Abs(mx)
	my = math.Abs(my)
	for i := range td.Data {
		td.Data[i][0] += mx
		td.Data[i][1] += my
		td.Data[i][0] *= 40
		td.Data[i][1] *= 40
		td.Data[i][0] += 50
		td.Data[i][1] += 50
	}

	image := gg.NewContext(size, size)
	image.DrawRectangle(0, 0, float64(size), float64(size))
	image.SetRGB(1, 1, 1)
	image.Fill()
	for i, p := range td.Data {
		if i != len(td.Data)-1 {
			image.DrawLine(p[0], p[1], td.Data[i+1][0], td.Data[i+1][1])
		}
	}
	image.SetRGB(0, 0, 0)
	image.Stroke()
	image.SavePNG("./test/example_track/track.png")

	return nil
}

func ExampleTrackObjParse() error {

	f, err := os.Open("./test/example_track/track.obj")
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewScanner(f)
	b.Split(bufio.ScanLines)

	t := []string{}

	for b.Scan() {
		t = append(t, b.Text())
	}

	ov := [][3]float64{}
	of := [][4]int{}

	for _, l := range t {
		str := strings.Split(l, " ")
		if len(str) > 0 {
			if str[0] == "v" {
				ovv := [3]float64{}
				for i, v := range str {
					if i > 0 {
						vv, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return err
						}
						ovv[i-1] = vv
					}
				}
				ov = append(ov, ovv)
			}
			if str[0] == "f" {
				ofv := [4]int{}
				for i, v := range str {
					if i > 0 {
						vv, err := strconv.Atoi(v)
						if err != nil {
							return err
						}
						ofv[i-1] = vv
					}
				}
				of = append(of, ofv)
			}
		}
	}

	ofo := [][4]int{}
	ofo = append(ofo, of[0])

	for len(ofo) < len(of) {
		for _, v := range of {
			if ofo[len(ofo)-1][2] == v[1] && ofo[len(ofo)-1][3] == v[0] {
				ofo = append(ofo, v)
				break
			}
		}
	}

	td := TrackData{}
	td.Data = make([][2]float64, 0)

	for _, v := range ofo {
		td.Data = append(td.Data, [2]float64{ov[v[0]-1][0], ov[v[0]-1][2]})
		td.Data = append(td.Data, [2]float64{ov[v[1]-1][0], ov[v[1]-1][2]})
	}

	for _, v := range td.Data {
		fmt.Printf("[%.2f, %.2f],", v[0], v[1])
	}

	return nil
}

type TrackData struct {
	Data [][2]float64 `json:"data"`
}

func index(v, max int) int {

	if v >= max {
		return v - max
	}
	return v

}

func distance(x0, y0, x1, y1 float64) float64 {

	return math.Sqrt(math.Pow(x1-x0, 2) + math.Pow(y1-y0, 2))
}

func deg(v float64) float64 {
	return v * (180.0 / math.Pi)
}

func angle(x0, y0, x1, y1, x2, y2 float64) float64 {
	x01 := x0 - x1
	x21 := x2 - x1
	y01 := y0 - y1
	y21 := y2 - y1
	return math.Acos((x01*x21 + y01*y21) / (math.Sqrt(x01*x01+y01*y01) * math.Sqrt(x21*x21+y21*y21)))
}

func lerp(x0, y0, x1, y1, t float64) (float64, float64) {
	x := x0 + (x1-x0)*t
	y := y0 + (y1-y0)*t
	return x, y
}

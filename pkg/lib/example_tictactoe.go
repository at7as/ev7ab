package lib

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/nsf/termbox-go"
)

var ExampleTictactoeDigits []rune = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

type ExampleTictactoe struct {
	table [][]rune
	step  int
	state []float64
	msg   string
	l     *lab.Lab
	id    int
}

func (p *ExampleTictactoe) Load(setup map[string]string) error {

	return nil
}

func (p *ExampleTictactoe) Setup(key, value string) error {

	return nil
}

func (p *ExampleTictactoe) Produce(n lab.Next, op lab.Next, r []float64) []float64 {

	in := make([]float64, 18)
	out := make([]float64, 9)
	excl := make([]int, 0)
	if len(r) == 0 {
		r = make([]float64, 3)
	}
	tie := false
	for i := range 5 {

		out = n(in)
		excl = make([]int, 0)
		for i := range 9 {
			if in[i] > 0.0 || in[i+9] > 0.0 {
				excl = append(excl, i)
			}
		}
		in[lab.IndexOfMaxExcluded(out, excl)] = 1.0
		if resultCheck(in[:9]) {
			r[0] += 1.0
			r[1] += 1.0
			break
		}

		if i != 4 {
			out = op(append(in[9:], in[:9]...))
			excl = make([]int, 0)
			for i := range 9 {
				if in[i] > 0.0 || in[i+9] > 0.0 {
					excl = append(excl, i)
				}
			}
			in[lab.IndexOfMaxExcluded(out, excl)+9] = 1.0
			if resultCheck(in[9:]) {
				r[0] += 1.0
				break
			}
		}

		if i == 4 {
			tie = true
		}

	}
	if tie {
		r[0] += 1.0
		r[2] += 1.0
	}

	return r
}

func (p *ExampleTictactoe) Validate(r []float64) bool {

	if r[0] > 0 && r[1] == 0.0 && r[2] == 0.0 {
		return false
	}

	return true
}

func (p *ExampleTictactoe) Compare(a, b []float64) bool {

	if b[1] == 0.0 {
		return true
	}
	if a[1] == 0.0 && b[1] != 0.0 {
		return false
	}

	return a[1]/a[0]*a[1] > b[1]/b[0]*b[1]
}

func (p *ExampleTictactoe) Best(v []float64) string {

	best := "---"
	if len(v) > 0 {
		best = fmt.Sprintf("%.1f", v[1]/v[0]*v[1])
	}

	return best
}

func (p *ExampleTictactoe) Goal(v []float64) bool {

	return false
}

func ExampleTictactoeApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleTictactoe{}, *cfgFile, true)

}

func resultValidate(in []float64, i int) bool {

	if in[i] > 0.0 || in[9+i] > 0.0 {
		return false
	}

	return true
}

func resultCheck(in []float64) bool {

	if (in[0] > 0.0 && in[1] > 0.0 && in[2] > 0.0) ||
		(in[3] > 0.0 && in[4] > 0.0 && in[5] > 0.0) ||
		(in[6] > 0.0 && in[7] > 0.0 && in[8] > 0.0) ||
		(in[0] > 0.0 && in[3] > 0.0 && in[6] > 0.0) ||
		(in[1] > 0.0 && in[4] > 0.0 && in[7] > 0.0) ||
		(in[2] > 0.0 && in[5] > 0.0 && in[8] > 0.0) ||
		(in[0] > 0.0 && in[4] > 0.0 && in[8] > 0.0) ||
		(in[2] > 0.0 && in[4] > 0.0 && in[6] > 0.0) {
		return true
	}

	return false
}

func ExampleTictactoeTry(id int) {

	f, err := os.Open("./test/example_tictactoe/ev.lab")
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

	example := &ExampleTictactoe{id: id}
	example.l = lab.New(example, false)
	if err = example.l.Import(b); err != nil {
		log.Panicln(err)
	}

	if err = termbox.Init(); err != nil {
		log.Panicln(err)
	}
	defer termbox.Close()

	example.reset()
	example.drawTable()

	run := true
	for run {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				run = false
			}
			if ev.Key == termbox.KeyEsc {
				example.reset()
				example.drawTable()
			}
			if slices.Contains(ExampleTictactoeDigits, ev.Ch) {
				next := false
				if example.step < 9 {
					next = example.playerMove(ev.Ch)
					if next {
						example.drawTable()
					}
				}
				if example.step < 9 && next {
					time.Sleep(time.Second)
					example.botMove()
					example.drawTable()
				}
				if example.step == 9 && example.msg == "" {
					example.msg = "Tie"
					example.drawTable()
				}
			}
		case termbox.EventError:
			log.Panicln(ev.Err)
		}

	}

}

func (p *ExampleTictactoe) reset() {

	p.state = make([]float64, 18)
	p.state[lab.IndexOfMax(p.l.ProjectValue(p.id, p.state))] = 1.0
	p.step = 1
	p.table = [][]rune{
		{'┌', '─', '─', '─', '┬', '─', '─', '─', '┬', '─', '─', '─', '┐'},
		{'1', ' ', ' ', ' ', '│', ' ', ' ', ' ', '│', ' ', ' ', ' ', '3'},
		{'├', '─', '─', '─', '┼', '─', '2', '─', '┼', '─', '─', '─', '┤'},
		{'4', ' ', ' ', ' ', '│', ' ', ' ', ' ', '│', ' ', ' ', ' ', '6'},
		{'├', '─', '─', '─', '┼', '─', '5', '─', '┼', '─', '─', '─', '┤'},
		{'7', ' ', ' ', ' ', '│', ' ', ' ', ' ', '│', ' ', ' ', ' ', '9'},
		{'└', '─', '─', '─', '┴', '─', '8', '─', '┴', '─', '─', '─', '┘'},
	}
	p.msg = ""

}

func (p *ExampleTictactoe) playerMove(char rune) bool {

	index, err := strconv.Atoi(string(char))
	if err != nil {
		log.Panicln(err)
	}
	index--

	if !resultValidate(p.state, index) {
		return false
	}

	p.state[9+index] = 1.0

	if resultCheck(p.state[9:]) {
		p.msg = "Player wins"
		p.step = 9
	} else {
		p.step++
	}

	return true
}

func (p *ExampleTictactoe) botMove() {

	excl := make([]int, 0)
	for i := range 9 {
		if p.state[i] > 0.0 || p.state[i+9] > 0.0 {
			excl = append(excl, i)
		}
	}

	p.state[lab.IndexOfMaxExcluded(p.l.ProjectValue(p.id, p.state), excl)] = 1.0

	if resultCheck(p.state[:9]) {
		p.msg = "Bot wins"
		p.step = 9
	} else {
		p.step++
	}

}

func (p *ExampleTictactoe) drawTable() {

	for i := range 9 {
		if p.state[i] > 0 {
			x, y := getPos(i)
			p.table[y*2+1][x*4+2] = 'X'
		}
	}

	for i := range 9 {
		if p.state[9+i] > 0 {
			x, y := getPos(i)
			p.table[y*2+1][x*4+2] = 'O'
		}
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y := range p.table {
		for x := range p.table[y] {
			termbox.SetChar(x+3, y+2, p.table[y][x])
		}
	}

	printMsg(3, 10, p.msg)

	termbox.Flush()

}

func getPos(i int) (int, int) {

	switch i {
	case 0:
		return 0, 0
	case 1:
		return 1, 0
	case 2:
		return 2, 0
	case 3:
		return 0, 1
	case 4:
		return 1, 1
	case 5:
		return 2, 1
	case 6:
		return 0, 2
	case 7:
		return 1, 2
	case 8:
		return 2, 2
	default:
		return 0, 0
	}

}

func printMsg(x, y int, msg string) {

	for _, c := range msg {
		termbox.SetChar(x, y, c)
		x++
	}

}

package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

// EditViewWidget ...
type EditViewWidget struct{}

// NewEditViewWidget ...
func NewEditViewWidget() *EditViewWidget {
	return &EditViewWidget{}
}

// Layout ...
func (w *EditViewWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("edit", -1, 0, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	v.Frame = true
	v.Title = " Edit "

	if app.edit != nil {

		caption := make([]string, maxX-7)
		for i := range caption {
			caption[i] = " "
		}
		bufCaption := make([]string, max((len(app.edit.n.model)-1)*5+3, maxX-7))

		footer := make([]string, maxX-7)
		for i := range footer {
			footer[i] = " "
		}
		bufFooter := make([]string, max((len(app.edit.n.model)-1)*5+3, maxX-7))

		buf := make([][]string, maxY-11)
		for i := range buf {
			buf[i] = make([]string, maxX-8)
			for ii := range buf[i] {
				buf[i][ii] = " "
			}
		}

		maxS := len(app.edit.n.model)
		maxN := 0
		for _, v := range app.edit.n.model {
			maxN = max(maxN, len(v.stage))
		}
		maxS--
		maxN += 2

		bufWhite := make([][]string, max(maxN*2, maxY-11+2))
		for i := range bufWhite {
			bufWhite[i] = make([]string, max(maxS*5+2, maxX-8))
		}

		bufGreen := make([][]string, max(maxN*2, maxY-11+2))
		for i := range bufGreen {
			bufGreen[i] = make([]string, max(maxS*5+2, maxX-8))
		}

		bufRed := make([][]string, max(maxN*2, maxY-11+2))
		for i := range bufRed {
			bufRed[i] = make([]string, max(maxS*5+2, maxX-8))
		}

		for i, s := range app.edit.n.model {

			i5 := i * 5

			fillStage := "7"
			if i == app.cursor.s {
				fillStage = "2"
			}
			if app.linkedit && i == app.link.s {
				fillStage = "2"
			}

			if i == 0 {
				bufCaption[i5] = fmt.Sprintf("\033[3%s;1mi", fillStage)
				bufCaption[i5+1] = fmt.Sprintf("n\033[0m")
			} else if i == len(app.edit.n.model)-1 {
				bufCaption[i5] = fmt.Sprintf("\033[3%s;1mo", fillStage)
				bufCaption[i5+1] = "u"
				bufCaption[i5+2] = fmt.Sprintf("t\033[0m")
			} else {
				bufCaption[i5] = fmt.Sprintf("\033[3%s;1m%v\033[0m", fillStage, i-1)
				if i > 10 {
					bufCaption[i5] = fmt.Sprintf("\033[3%s;1m%s", fillStage, strconv.Itoa(i - 1)[0:1])
					bufCaption[i5+1] = fmt.Sprintf("%s\033[0m", strconv.Itoa(i - 1)[1:2])
				}
				if i > 100 {
					bufCaption[i5] = fmt.Sprintf("\033[3%s;1m%s", fillStage, strconv.Itoa(i - 1)[0:1])
					bufCaption[i5+1] = strconv.Itoa(i - 1)[1:2]
					bufCaption[i5+2] = fmt.Sprintf("%s\033[0m", strconv.Itoa(i - 1)[2:3])
				}
			}

			if i == app.cursor.s {
				if len(app.edit.n.model[i].stage) > 0 {
					bufFooter[i5] = fmt.Sprintf("\033[32;1m%v\033[0m", app.cursor.n)
					if app.cursor.n > 9 {
						bufFooter[i5] = fmt.Sprintf("\033[32;1m%s", strconv.Itoa(app.cursor.n)[0:1])
						bufFooter[i5+1] = fmt.Sprintf("%s\033[0m", strconv.Itoa(app.cursor.n)[1:2])
					}
					if app.cursor.n > 99 {
						bufFooter[i5] = fmt.Sprintf("\033[32;1m%s", strconv.Itoa(app.cursor.n)[0:1])
						bufFooter[i5+1] = strconv.Itoa(app.cursor.n)[1:2]
						bufFooter[i5+2] = fmt.Sprintf("%s\033[0m", strconv.Itoa(app.cursor.n)[2:3])
					}
				} else {
					bufFooter[i5] = fmt.Sprintf("\033[32;1m%s", "─")
					bufFooter[i5+1] = fmt.Sprintf("%s\033[0m", "─")
				}
			}

			for ii, n := range s.stage {

				oneii2 := 1 + ii*2

				fillNode := "7"
				fillNodeChar := "█"
				if i == app.cursor.s && ii == app.cursor.n {
					fillNode = "5"
				}
				if app.linkedit && i == app.link.s && ii == app.link.n {
					fillNodeChar = "▓"
				}
				bufWhite[oneii2][i5] = fmt.Sprintf("\033[3%s;4m%s", fillNode, fillNodeChar)
				bufWhite[oneii2][i5+1] = fmt.Sprintf("%s\033[0m", fillNodeChar)

				if len(n.source) > 0 {
					c := "7"
					if i == app.cursor.s && ii == app.cursor.n {
						c = "2"
					}
					bufWhite[oneii2][i5-1] = fmt.Sprintf("\033[3%s;4m─\033[0m", c)
				}

				for _, src := range n.source {

					bufWhite[1+src.n*2][src.s*5+2] = "─"

					s53 := i*5 - 2
					srcn := src.n

					if src.s != i-1 {

						srcn = len(app.edit.n.model[i-1].stage)
						bufWhite[1+srcn*2][i*5-3] = "╾"

						iin := len(app.edit.n.model[src.s+1].stage)
						x := src.s*5 + 3
						bufWhite[1+iin*2][x+1] = "╼"

						if src.n < iin {
							bufWhite[1+src.n*2][x] = drawBox(bufWhite[1+src.n*2][x], "┐")
							bufWhite[1+iin*2][x] = drawBox(bufWhite[1+iin*2][x], "└")
							for y := 1 + src.n*2 + 1; y < 1+iin*2; y++ {
								bufWhite[y][x] = drawBox(bufWhite[y][x], "│")
							}
						} else if src.n > iin {
							bufWhite[1+src.n*2][x] = drawBox(bufWhite[1+src.n*2][x], "┘")
							bufWhite[1+iin*2][x] = drawBox(bufWhite[1+iin*2][x], "┌")
							for y := 1 + iin*2 + 1; y < 1+src.n*2; y++ {
								bufWhite[y][x] = drawBox(bufWhite[y][x], "│")
							}
						} else {
							bufWhite[1+src.n*2][x] = drawBox(bufWhite[1+src.n*2][x], "─")
						}

					}

					if srcn < ii {
						bufWhite[1+srcn*2][s53] = drawBox(bufWhite[1+srcn*2][s53], "┐")
						bufWhite[1+ii*2][s53] = drawBox(bufWhite[1+ii*2][s53], "└")
						for y := 1 + srcn*2 + 1; y < 1+ii*2; y++ {
							bufWhite[y][s53] = drawBox(bufWhite[y][s53], "│")
						}
					} else if srcn > ii {
						bufWhite[1+srcn*2][s53] = drawBox(bufWhite[1+srcn*2][s53], "┘")
						bufWhite[1+ii*2][s53] = drawBox(bufWhite[1+ii*2][s53], "┌")
						for y := 1 + ii*2 + 1; y < 1+srcn*2; y++ {
							bufWhite[y][s53] = drawBox(bufWhite[y][s53], "│")
						}
					} else {
						bufWhite[1+srcn*2][s53] = drawBox(bufWhite[1+srcn*2][s53], "─")
					}

				}

			}
		}

		if len(app.edit.n.model[app.cursor.s].stage) > 0 {
			for _, src := range app.edit.n.model[app.cursor.s].stage[app.cursor.n].source {
				fillNodeChar := "█"
				if app.linkedit && src.s == app.link.s && src.n == app.link.n {
					fillNodeChar = "▓"
				}
				bufGreen[1+src.n*2][src.s*5] = fmt.Sprintf("%s", fillNodeChar)
				bufGreen[1+src.n*2][src.s*5+1] = fmt.Sprintf("%s", fillNodeChar)
				bufGreen[1+src.n*2][src.s*5+2] = "─"

				i := app.cursor.s
				ii := app.cursor.n

				s53 := i*5 - 2
				srcn := src.n

				if src.s != i-1 {

					srcn = len(app.edit.n.model[i-1].stage)
					bufGreen[1+srcn*2][i*5-3] = "╾"

					iin := len(app.edit.n.model[src.s+1].stage)
					x := src.s*5 + 3
					bufGreen[1+iin*2][x+1] = "╼"

					if src.n < iin {
						bufGreen[1+src.n*2][x] = drawBox(bufGreen[1+src.n*2][x], "┐")
						bufGreen[1+iin*2][x] = drawBox(bufGreen[1+iin*2][x], "└")
						for y := 1 + src.n*2 + 1; y < 1+iin*2; y++ {
							bufGreen[y][x] = drawBox(bufGreen[y][x], "│")
						}
					} else if src.n > iin {
						bufGreen[1+src.n*2][x] = drawBox(bufGreen[1+src.n*2][x], "┘")
						bufGreen[1+iin*2][x] = drawBox(bufGreen[1+iin*2][x], "┌")
						for y := 1 + iin*2 + 1; y < 1+src.n*2; y++ {
							bufGreen[y][x] = drawBox(bufGreen[y][x], "│")
						}
					} else {
						bufGreen[1+src.n*2][x] = drawBox(bufGreen[1+src.n*2][x], "─")
					}

				}

				if srcn < ii {
					bufGreen[1+srcn*2][s53] = drawBox(bufGreen[1+srcn*2][s53], "┐")
					bufGreen[1+ii*2][s53] = drawBox(bufGreen[1+ii*2][s53], "└")
					for y := 1 + srcn*2 + 1; y < 1+ii*2; y++ {
						bufGreen[y][s53] = drawBox(bufGreen[y][s53], "│")
					}
				} else if srcn > ii {
					bufGreen[1+srcn*2][s53] = drawBox(bufGreen[1+srcn*2][s53], "┘")
					bufGreen[1+ii*2][s53] = drawBox(bufGreen[1+ii*2][s53], "┌")
					for y := 1 + ii*2 + 1; y < 1+srcn*2; y++ {
						bufGreen[y][s53] = drawBox(bufGreen[y][s53], "│")
					}
				} else {
					bufGreen[1+srcn*2][s53] = drawBox(bufGreen[1+srcn*2][s53], "─")
				}

			}
		}

		for y := range bufGreen {
			for x := range bufGreen[y] {
				if bufGreen[y][x] != "" {
					bufGreen[y][x] = fmt.Sprintf("\033[32;4m%s\033[0m", bufGreen[y][x])
				}
			}
		}

		if app.invalid {
			for i, s := range app.edit.n.model {
				if len(s.stage) == 0 {
					for iii := range bufRed {
						bufRed[iii][i*5] = fmt.Sprintf("\033[31;1m█")
						bufRed[iii][i*5+1] = fmt.Sprintf("█\033[0m")
					}
				}
				for ii, n := range s.stage {
					if !n.valid {
						bufRed[1+ii*2][i*5] = fmt.Sprintf("\033[31;1m█")
						bufRed[1+ii*2][i*5+1] = fmt.Sprintf("█\033[0m")
					}
				}
			}
		}

		maxXs := int(math.Floor(float64(maxX)-8.0) / 5.0)
		maxYn := int(math.Floor(float64(maxY)-11.0)/2.0) - 2
		if app.linkedit {
			if app.link.s > app.scroll.s+maxXs {
				app.scroll.s = app.link.s - maxXs
			}
			if app.link.s < app.scroll.s {
				app.scroll.s = app.link.s
			}
			if app.link.n > app.scroll.n+maxYn {
				app.scroll.n = app.link.n - maxYn
			}
			if app.link.n < app.scroll.n {
				app.scroll.n = app.link.n
			}
		} else {
			if app.cursor.s > app.scroll.s+maxXs {
				app.scroll.s = app.cursor.s - maxXs
			}
			if app.cursor.s < app.scroll.s {
				app.scroll.s = app.cursor.s
			}
			if app.cursor.n > app.scroll.n+maxYn {
				app.scroll.n = app.cursor.n - maxYn
			}
			if app.cursor.n < app.scroll.n {
				app.scroll.n = app.cursor.n
			}
		}

		for i := range buf {
			for ii := range buf[i] {
				y := i + app.scroll.n*2
				x := ii + app.scroll.s*5
				if bufWhite[y][x] != "" {
					buf[i][ii] = bufWhite[y][x]
				}
				if bufGreen[y][x] != "" {
					buf[i][ii] = bufGreen[y][x]
				}
				if bufRed[y][x] != "" {
					buf[i][ii] = bufRed[y][x]
				}
			}
		}

		for i := range caption {
			if bufCaption[i+app.scroll.s*5] != "" {
				caption[i] = bufCaption[i+app.scroll.s*5]
			}
		}

		for i := range footer {
			if bufFooter[i+app.scroll.s*5] != "" {
				footer[i] = bufFooter[i+app.scroll.s*5]
			}
		}

		fmt.Fprintf(v, " \n")
		fmt.Fprintf(v, "  %s %s      %s\n", dirtyMark(app.edit.n.dirty), untrimRight(app.edit.status.text(), 12), numberMark(app.edit))
		fmt.Fprintf(v, " \n")
		nodeStat := ""
		if len(app.edit.n.model[app.cursor.s].stage) > 0 {
			nodeStat = fmt.Sprintf("node = %v : %v", app.edit.n.model[app.cursor.s].stage[app.cursor.n].size, app.edit.n.model[app.cursor.s].stage[app.cursor.n].volume)
		}
		fmt.Fprintf(v, "    total = %v : %v    stage = %v : %v    %s\n", app.edit.n.size, app.edit.n.volume, app.edit.n.model[app.cursor.s].size, app.edit.n.model[app.cursor.s].volume, nodeStat)
		fmt.Fprintf(v, " \n")
		fmt.Fprintf(v, "  ┌ %s┐  \n", strings.Join(caption, ""))
		for _, s := range buf {
			fmt.Fprintf(v, "  │ %s │ \n", strings.Join(s, ""))
		}
		fmt.Fprintf(v, "  └ %s┘  \n", strings.Join(footer, ""))
	}

	return nil
}

var charBox map[string]int = map[string]int{
	" ": 0b0000,
	"└": 0b1100,
	"┌": 0b0110,
	"┐": 0b0011,
	"┘": 0b1001,
	"┴": 0b1101,
	"├": 0b1110,
	"┬": 0b0111,
	"┤": 0b1011,
	"│": 0b1010,
	"─": 0b0101,
	"┼": 0b1111,
}

var charInt map[int]string = map[int]string{
	0b0000: " ",
	0b1100: "└",
	0b0110: "┌",
	0b0011: "┐",
	0b1001: "┘",
	0b1101: "┴",
	0b1110: "├",
	0b0111: "┬",
	0b1011: "┤",
	0b1010: "│",
	0b0101: "─",
	0b1111: "┼",
}

func drawBox(src, char string) string {

	if src == "" {
		src = " "
	}

	return charInt[charBox[src]|charBox[char]]
}

func dirtyMark(d bool) string {

	if d {
		return "\033[32;1m*\033[0m"
	}

	return " "
}

func numberMark(p *project) string {

	if p.status == psNew {

		if p.o != nil {
			return fmt.Sprintf("## ? \033[30;1m[%v]\033[0m", p.o.id)
		}
		return fmt.Sprintf("## ?")

	}

	if p.o != nil {
		return fmt.Sprintf("## %v \033[30;1m[%v]\033[0m", p.id, p.o.id)
	}
	return fmt.Sprintf("## %v", p.id)

}

func selectStageLeft() {

	if app.linkedit {
		if app.link.s > 0 {
			app.link.s--
		}
	} else {
		if app.cursor.s > 0 {
			app.cursor.s--
		}
	}

}

func selectStageRight() {

	if app.linkedit {
		if app.link.s < len(app.edit.n.model)-1 {
			app.link.s++
		}
	} else {
		if app.cursor.s < len(app.edit.n.model)-1 {
			app.cursor.s++
		}
	}

}

func selectNodeUp() {

	if app.linkedit {
		if app.link.n > 0 {
			app.link.n--
		}
	} else {
		if app.cursor.n > 0 {
			app.cursor.n--
		}
	}

}

func selectNodeDown() {

	if app.linkedit {
		if app.link.n < len(app.edit.n.model[app.link.s].stage)-1 {
			app.link.n++
		}
	} else {
		if app.cursor.n < len(app.edit.n.model[app.cursor.s].stage)-1 {
			app.cursor.n++
		}
	}

}

func insertStage() {

	selectStageRight()
	app.edit.n.addStage(app.cursor.s)

}

func deleteStage() {

	app.edit.n.removeStage(app.cursor.s)

}

func shiftCursorY() {

	if app.linkedit {
		if app.link.n != 0 && app.link.n > len(app.edit.n.model[app.link.s].stage)-1 {
			app.link.n = max(0, len(app.edit.n.model[app.link.s].stage)-1)
		}
	} else {
		if app.cursor.n != 0 && app.cursor.n > len(app.edit.n.model[app.cursor.s].stage)-1 {
			app.cursor.n = max(0, len(app.edit.n.model[app.cursor.s].stage)-1)
		}
	}

}

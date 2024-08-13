package app

import (
	"fmt"
	"math"
	"strings"

	"github.com/jroimartin/gocui"
)

var width map[string]int

func init() {

	width = make(map[string]int)
	width["num"] = 0
	width["status"] = 0
	width["size"] = 0
	width["volume"] = 0
	width["gen"] = 0
	width["ev"] = 0
	width["age"] = 0
	width["best"] = 0

}

type resultViewWidget struct{}

func newResultViewWidget() *resultViewWidget {
	return &resultViewWidget{}
}

func (w *resultViewWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("result", -1, 0, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	v.Frame = true
	v.Title = " Result "

	// top to result when project is goal

	for k := range width {
		width[k] = 0
	}

	for _, p := range app.result.list {
		width["num"] = max(len(fmt.Sprint(p.id)), width["num"])
		width["status"] = max(len(p.status.text()), width["status"])
		width["size"] = max(len(fmt.Sprint(p.s.size)), width["size"])
		width["volume"] = max(len(fmt.Sprint(p.s.volume)), width["volume"])
		width["gen"] = max(len(fmt.Sprint(p.s.gen)), width["gen"])
		width["ev"] = max(len(fmt.Sprint(p.s.ev)), width["ev"])
		width["age"] = max(len(fmt.Sprint(p.s.age)), width["age"])
		width["best"] = max(len(p.s.best), width["best"])
	}

	width["num"] = max(width["num"]+2, 7)
	width["status"] = max(width["status"]+2, 10)
	width["size"] = max(width["size"]+2, 9)
	width["volume"] = max(width["volume"]+2, 11)
	width["gen"] = max(width["gen"]+2, 8)
	width["ev"] = max(width["ev"]+2, 8)
	width["age"] = max(width["age"]+2, 8)
	width["best"] = max(width["best"]+2, 9)

	fmt.Fprintf(v, " \n")
	fmt.Fprintf(v, "  %s│", space("##", 1, width["num"]-3))
	if app.result.holded {
		fmt.Fprintf(v, "%s│", space("\033[33;4mSTATUS\033[0m", 1, width["status"]-7))
	} else {
		fmt.Fprintf(v, "%s│", space("STATUS", 1, width["status"]-7))
	}
	fmt.Fprintf(v, "%s│", space("SIZE", width["size"]-5, 1))
	fmt.Fprintf(v, "%s│", space("VOLUME", width["volume"]-7, 1))
	fmt.Fprintf(v, "%s│", space("GEN", width["gen"]-4, 1))
	fmt.Fprintf(v, "%s│", space("EV", width["ev"]-3, 1))
	fmt.Fprintf(v, "%s│", space("AGE", width["age"]-4, 1))
	fmt.Fprintf(v, "%s\n", space("BEST", width["best"]-5, 1))

	buf := make([]string, maxY-5)
	for i := range buf {
		buf[i] = " \n"
	}

	bufList := make([]string, 0, max(maxY-5, len(app.result.list)*2+1))

	app.result.update()

	for i, p := range app.result.show {

		bufList = append(bufList, fmt.Sprintf("  %s┼%s┼%s┼%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", width["num"]), strings.Repeat("─", width["status"]), strings.Repeat("─", width["size"]), strings.Repeat("─", width["volume"]), strings.Repeat("─", width["gen"]), strings.Repeat("─", width["ev"]), strings.Repeat("─", width["age"]), strings.Repeat("─", width["best"])))

		selected := " "
		if p.sel {
			selected = "█"
		}

		fontColor := ""
		fontStatusColor := ""
		if i == app.position && !app.run {
			if p.s.goal {
				fontColor = "2;7"
				fontStatusColor = "2;7"
			} else {
				fontColor = "7;7"
				fontStatusColor = "7;7"
			}
		} else {
			if p.status == psActive {
				fontColor = "7;4"
				if p.s.goal {
					fontStatusColor = "2;1"
				} else {
					fontStatusColor = "7;4"
				}
			} else if p.status == psHolded {
				fontColor = "0;1"
				fontStatusColor = "3;4"
			} else {
				fontColor = "0;1"
				fontStatusColor = "1;1"
			}
			if p.s.goal {
				fontColor = "2;1"
			}
		}

		bufList = append(bufList, fmt.Sprintf("%s \033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m\n", selected, fontColor, untrimRight(fmt.Sprint(p.id), width["num"]), fontStatusColor, untrimRight(p.status.text(), width["status"]), fontColor, untrimLeft(fmt.Sprint(p.s.size), width["size"]), fontColor, untrimLeft(fmt.Sprint(p.s.volume), width["volume"]), fontColor, untrimLeft(fmt.Sprint(p.s.gen), width["gen"]), fontColor, untrimLeft(fmt.Sprint(p.s.ev), width["ev"]), fontColor, untrimLeft(fmt.Sprint(p.s.age), width["age"]), fontColor, untrimLeft(p.s.best, width["best"])))

	}

	bufList = append(bufList, fmt.Sprintf("  %s┴%s┴%s┴%s┴%s┴%s┴%s┴%s\n", strings.Repeat("─", width["num"]), strings.Repeat("─", width["status"]), strings.Repeat("─", width["size"]), strings.Repeat("─", width["volume"]), strings.Repeat("─", width["gen"]), strings.Repeat("─", width["ev"]), strings.Repeat("─", width["age"]), strings.Repeat("─", width["best"])))

	cnt := int(max(maxY-5, len(app.result.list)*2+1)) - len(bufList)
	for range cnt {
		bufList = append(bufList, "")
	}

	maxP := int(math.Floor(float64(maxY)-5.0) / 2.0)
	if app.position > app.offset+maxP {
		app.offset = app.position - maxP
	}
	if app.position < app.offset {
		app.offset = app.position
	}

	for i := range buf {
		y := i + app.offset*2
		if bufList[y] != "" {
			buf[i] = bufList[y]
		}
	}

	for _, s := range buf {
		fmt.Fprint(v, s)
	}

	return nil
}

package app

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

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

	// 80

	// 80 for

	fmt.Fprintf(v, "  %s  %s│\n", space(5), space(5))
	fmt.Fprintf(v, "  %s##%s│", space(5), space(5))
	fmt.Fprintf(v, "%sSIZE%s│", space(4), space(4))
	fmt.Fprintf(v, "%sVOLUME%s│", space(3), space(3))
	fmt.Fprintf(v, "%sGEN%s│", space(4), space(5))
	fmt.Fprintf(v, "%sEV%s│", space(5), space(5))
	fmt.Fprintf(v, "%sAGE%s│", space(4), space(5))
	fmt.Fprintf(v, "%sBEST%s│", space(4), space(4))
	if app.result.holded {
		fmt.Fprintf(v, "%s\033[33;4mSTATUS\033[0m%s\n", space(3), space(3))
	} else {
		fmt.Fprintf(v, "%sSTATUS%s\n", space(3), space(3))
	}

	cnt := 0
	for i, p := range app.result.l {
		if (!app.result.holded && p.status == psActive) || (app.result.holded && p.status != psTerminated) {
			cnt++
			fmt.Fprintf(v, "  %s│%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12))

			selected := " "
			if p.sel {
				selected = "█"
			}
			if i == app.result.ai && !app.run {
				fmt.Fprintf(v, "%s \033[37;7m %v\033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m %v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

			} else {
				if p.status == psActive {
					fmt.Fprintf(v, "%s  %v│%v │%v │%v │%v │ %v\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				} else if p.status == psHolded {
					fmt.Fprintf(v, "%s  \033[30;1m%v\033[0m│\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │ \033[33;4m%v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				} else {
					fmt.Fprintf(v, "%s  \033[30;1m%v\033[0m│\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │ \033[31;1m%v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				}
			}
		}
	}

	if cnt == 0 {
		fmt.Fprintf(v, "  %s│%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12))
	}

	for range maxY - cnt*2 - 5 {
		fmt.Fprintf(v, "  %s  %s│\n", space(5), space(5))
	}

	return nil
}

func space(n int) string {

	return strings.Repeat(" ", n)
}

func untrimRight(v string, n int) string {

	return v + strings.Repeat(" ", n-len(v))
}

func untrimLeft(v string, n int) string {

	return strings.Repeat(" ", n-len(v)) + v
}

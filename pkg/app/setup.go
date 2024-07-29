package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type setupViewWidget struct{}

func newSetupViewWidget() *setupViewWidget {
	return &setupViewWidget{}
}

func (w *setupViewWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("setup", -1, 0, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	v.Frame = true
	v.Title = " Setup "

	fmt.Fprintf(v, " \n")
	for i, it := range app.setup.l {
		if i == app.setup.ai && !app.run {
			fmt.Fprintf(v, "  \033[37;7m%s  %s\033[0m\n\n", it.title, it.value)
		} else {
			fmt.Fprintf(v, "  %s  %s\n\n", it.title, it.value)
		}
	}

	return nil
}

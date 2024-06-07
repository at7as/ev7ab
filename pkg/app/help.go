package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// HelpWindowWidget ...
type HelpWindowWidget struct{}

// NewHelpWindowWidget ...
func NewHelpWindowWidget() *HelpWindowWidget {
	return &HelpWindowWidget{}
}

// Layout ...
func (w *HelpWindowWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("help", 6, 3, maxX-7, maxY-5)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()

	v.Frame = true
	v.Title = " Help "
	fmt.Fprintf(v, "%v\n", app.view)
	for i := 0; i <= 7; i++ {
		for _, j := range []int{1, 4, 7} {
			fmt.Fprintf(v, "\033[3%d;%dm#\033[0m", i, j)
		}
	}

	return nil
}

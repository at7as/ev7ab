package app

import "github.com/jroimartin/gocui"

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

	if len(app.result.l) == 0 {
		// make new
		// edit new
	}

	if len(app.result.l) > 0 {
		// edit current
	}

	// if make copy
	// edit copy

	// buffer hold until ESC
	// if ESC
	// make new
	// show current of project status

	return nil
}

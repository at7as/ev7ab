package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type widget struct {
	ctrl           controller
	name           string
	view           *gocui.View
	dirty          bool
	buf            []string
	x, y           int
	cursor, offset position
}

func (w *widget) setCursor(v position) {

	if v.x != w.cursor.x || v.y != w.cursor.y {
		w.cursor = v
		w.mark()
	}

}

func (w *widget) setOffset(v position) {

	if v.x != w.offset.x || v.y != w.offset.y {
		w.offset = v
		w.mark()
	}

}

func newWidget(ctrl controller, name string) *widget {

	return &widget{
		ctrl:   ctrl,
		name:   name,
		dirty:  true,
		buf:    []string{},
		x:      0,
		y:      0,
		cursor: newPosition(0, 0),
		offset: newPosition(0, 0),
	}
}

func (w *widget) Layout(g *gocui.Gui) error {

	if err = w.draw(); err != nil {
		return err
	}
	if w.dirty {
		if w.buf, err = w.ctrl.render(); err != nil {
			return err
		}
		w.dirty = false
	}

	w.view.Clear()
	for i, line := range w.buf {
		if i != len(w.buf)-1 {
			fmt.Fprintf(w.view, "%s\n", line)
		} else {
			fmt.Fprint(w.view, line)
		}
	}

	return nil
}

func (w *widget) draw() error {

	x, y := gui.Size()
	if x != w.x || y != w.y {
		w.x = x
		w.y = y
		w.mark()
	}

	x0, y0, x1, y1 := w.ctrl.transform(x, y)
	if w.view, err = gui.SetView(w.name, x0, y0, x1, y1); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func (w *widget) keybinding() error {

	return w.ctrl.keybinding()
}

func (w *widget) mark() {

	w.dirty = true

}

func (w *widget) clean() error {

	if w.view != nil {
		gui.DeleteKeybindings(w.view.Name())
		if err = gui.DeleteView(w.view.Name()); err != nil && err != gocui.ErrUnknownView {
			return err
		}
	}
	w.view = nil
	gui.Cursor = false

	return nil
}

func (w *widget) reset() {

	w.setCursor(newPosition(0, 0))
	w.setOffset(newPosition(0, 0))

}

type controller interface {
	transform(x int, y int) (int, int, int, int)
	render() ([]string, error)
	keybinding() error
}

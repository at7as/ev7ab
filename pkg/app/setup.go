package app

import (
	"fmt"
	"math"

	"github.com/jroimartin/gocui"
)

type setupWidget struct {
	*widget
	list           []kv
	width          int
	cursor, offset position
}

func (w *setupWidget) setList(key, value string) {

	for i, v := range w.list {
		if v.key == key {
			w.list[i].value = value
		}
	}
	w.dig()

}

func (w *setupWidget) setCursor(v position) {

	if v.x != w.cursor.x || v.y != w.cursor.y {
		w.cursor.x = v.x
		w.cursor.y = v.y
		w.dig()
	}

}

func (w *setupWidget) getCursor() position {

	return w.cursor
}

func (w *setupWidget) setOffset(v position) {

	if v.x != w.offset.x || v.y != w.offset.y {
		w.offset.x = v.x
		w.offset.y = v.y
		w.dig()
	}

}

func (w *setupWidget) getOffset() position {

	return w.offset
}

func newSetupWidget() *setupWidget {

	w := &setupWidget{
		list: []kv{
			{"LabFile", "./ev.lab"},
			{"Size", "1000"},
			{"Aggr", "avg"},
			{"Proc", "linear"},
			{"Goal", "false"},
			{"Duel", "false"},
			{"In", "2"},
			{"Out", "2"},
			{"InputFile", ""},
			{"Target", ""},
			{"Limit", ""},
		},
		width:  0,
		cursor: position{0, 0},
		offset: position{0, 0},
	}
	w.widget = newWidget(w, "setup")

	for _, v := range w.list {
		w.width = max(len(v.key), w.width)
	}

	return w
}

func (w *setupWidget) transform(x int, y int) (int, int, int, int) {

	return -1, 0, x, y - 1
}

func (w *setupWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, y := gui.Size()

	buf := make([]string, y-2)

	height := int(math.Floor((float64(y) - 4.0) / 2.0))
	width := x - w.width - 4

	if w.getCursor().y > w.getOffset().y+height-1 {
		w.setOffset(position{0, w.getCursor().y - height + 1})
	}
	if w.getCursor().y < w.getOffset().y {
		w.setOffset(position{0, w.getCursor().y})
	}

	buf[0] = " "
	for i := range min(len(w.list), height) {
		item := w.list[w.getOffset().y+i]
		buf[1+i*2+0] = fmt.Sprintf(" \033[3%sm%s  %s\033[0m ", w.focused(w.getOffset().y+i), space(item.key, 0, w.width-len(item.key)), item.value[:min(len(item.value), width)])
		buf[1+i*2+1] = " "
	}

	if w.getOffset().y > 0 {
		buf[0] = space("↑", x-1, 0)
	}

	if len(w.list)-w.getOffset().y > height {
		buf[y-3] = space("↓", x-1, 0)
	}

	return buf, nil
}

func (w *setupWidget) keybinding() error {

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.moveDown); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.moveUp); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.openSetupItem); err != nil {
		return err
	}

	return nil
}

func (w *setupWidget) moveDown(_ *gocui.Gui, _ *gocui.View) error {

	if app.status == appIdle {
		if w.getCursor().y < len(w.list)-1 {
			w.setCursor(position{0, w.getCursor().y + 1})
		}
	} else {

		_, y := gui.Size()
		height := int(math.Floor(float64(y)-4.0) / 2.0)
		yy := w.getOffset().y + height

		if yy >= len(w.list)-1 {
			w.setCursor(position{0, len(w.list) - 1})
		} else {
			w.setCursor(position{0, yy})
			w.setOffset(position{0, w.getOffset().y + 1})
		}

	}

	return nil
}

func (w *setupWidget) moveUp(_ *gocui.Gui, _ *gocui.View) error {

	if app.status == appIdle {
		if w.getCursor().y > 0 {
			w.setCursor(position{0, w.getCursor().y - 1})
		}
	} else {
		if w.getOffset().y > 0 {
			w.setCursor(position{0, w.getOffset().y - 1})
			w.setOffset(position{0, w.getOffset().y})
		}
	}

	return nil
}

func (w *setupWidget) openSetupItem(_ *gocui.Gui, _ *gocui.View) error {

	return app.openModal(newSetupItemBox(w.list[w.cursor.y].key, w.list[w.cursor.y].value))
}

func (w *setupWidget) focused(y int) string {

	if w.getCursor().y == y && app.status == appIdle {
		return "7;7"
	}

	return "7;4"
}

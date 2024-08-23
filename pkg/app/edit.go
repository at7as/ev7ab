package app

import (
	"fmt"
	"math"
	"time"

	"github.com/jroimartin/gocui"
)

type editWidget struct {
	*widget
	draft        *project
	mode         editMode
	invalid      bool
	link, scroll position
}

type editMode int

const (
	editNode editMode = iota
	editLink
)

func (w *editWidget) setDraft(p *project) {

	if p != w.draft {
		w.draft = p
		w.mark()
	}

}

func (w *editWidget) setMode(m editMode) {

	if m != w.mode {
		w.mode = m
		w.mark()
	}

}

func (w *editWidget) setInvalid(i bool) {

	if i != w.invalid {
		w.invalid = i
		w.mark()
	}

}

func (w *editWidget) setLink(v position) {

	if v.x != w.link.x || v.y != w.link.y {
		w.link.x = v.x
		w.link.y = v.y
		w.mark()
	}

}

func (w *editWidget) setScroll(v position) {

	if v.x != w.scroll.x || v.y != w.scroll.y {
		w.scroll.x = v.x
		w.scroll.y = v.y
		w.mark()
	}

}

func (w *editWidget) getX() int {

	if w.links() {
		return w.scroll.x
	}
	return w.offset.x
}

func (w *editWidget) getY() int {

	if w.links() {
		return w.scroll.y
	}
	return w.offset.y
}

func (w *editWidget) links() bool {

	return w.mode == editLink
}

func (w *editWidget) setNodeSize(v int) {

	if v != w.draft.n.model[w.cursor.x].stage[w.cursor.y].size {
		w.draft.n.model[w.cursor.x].stage[w.cursor.y].size = v
		w.mark()
		w.draft.n.measure()
	}

}

func (w *editWidget) getNodeSize() int {

	return w.draft.n.model[w.cursor.x].stage[w.cursor.y].size
}

func newEditWidget() *editWidget {

	w := &editWidget{
		draft:   nil,
		mode:    editNode,
		invalid: false,
		link:    position{0, 0},
		scroll:  position{0, 0},
	}
	w.widget = newWidget(w, "edit")

	return w
}

func (w *editWidget) transform(x int, y int) (int, int, int, int) {

	return -1, 0, x, y - 1
}

func (w *editWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, y := gui.Size()

	buf := make([]string, 5)

	if w.draft != nil {

		buf[0] = " "
		buf[1] = w.headerTop()
		buf[2] = " "
		buf[3] = w.headerBottom()
		buf[4] = " "

		buf = w.body(buf, x, y)

	}

	return buf, nil
}

func (w *editWidget) keybinding() error {

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlN, gocui.ModNone, w.newProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlE, gocui.ModNone, w.editProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlD, gocui.ModNone, w.dubProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlS, gocui.ModNone, w.saveProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlV, gocui.ModNone, w.validateProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.keyArrowLeft); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.keyArrowRight); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.keyArrowUp); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.keyArrowDown); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyInsert, gocui.ModNone, w.insertStage); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyDelete, gocui.ModNone, w.deleteStage); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.insertNode); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.editNodeSource); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlSpace, gocui.ModNone, w.editNode); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyBackspace, gocui.ModNone, w.deleteNode); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.linkEsc); err != nil {
		return err
	}

	return nil
}

func (w *editWidget) newProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)
	w.setDraft(newProject(nil))
	w.reset()

	return nil
}

func (w *editWidget) editProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	if w.draft != nil {
		w.draft.edit()
	}

	return nil
}

func (w *editWidget) dubProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	p := app.v.result.focused()
	if p != nil {
		w.setDraft(newProject(p))
		w.reset()
	}

	return nil
}

func (w *editWidget) saveProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	if w.draft != nil {
		w.draft.save()
	}

	return nil
}

func (w *editWidget) validateProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	if w.draft != nil {
		if !w.draft.validate() {
			w.setInvalid(true)
			go w.hideInvalid()
		}
	}

	return nil
}

func (w *editWidget) keyArrowLeft(_ *gocui.Gui, _ *gocui.View) error {

	w.selectStageLeft()
	w.shiftCursorY()

	return nil
}

func (w *editWidget) keyArrowRight(_ *gocui.Gui, _ *gocui.View) error {

	w.selectStageRight()
	w.shiftCursorY()

	return nil
}

func (w *editWidget) keyArrowUp(_ *gocui.Gui, _ *gocui.View) error {

	w.selectNodeUp()

	return nil
}

func (w *editWidget) keyArrowDown(_ *gocui.Gui, _ *gocui.View) error {

	w.selectNodeDown()

	return nil
}

func (w *editWidget) insertStage(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	w.setMode(editNode)
	insertStage()
	w.shiftCursorY()

	return nil
}

func (w *editWidget) deleteStage(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	w.setMode(editNode)
	deleteStage()
	w.shiftCursorY()
	w.draft.n.measure()

	return nil
}

func (w *editWidget) insertNode(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	w.setMode(editNode)

	if w.cursor.x > 0 && w.cursor.x < len(w.draft.n.model)-1 {
		w.draft.n.model[w.cursor.x].addNode()
		w.setCursor(position{len(w.draft.n.model[w.cursor.x].stage) - 1, w.cursor.y})
		return app.openModal(newNodeSizeBox())
	}

	w.draft.n.measure()

	return nil
}

func (w *editWidget) editNodeSource(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	if w.links() {
		if w.cursor.x != w.link.x || w.cursor.y != w.link.y {

			found := false
			for i, src := range w.draft.n.model[w.cursor.x].stage[w.cursor.y].source {
				if src.s == w.link.x && src.n == w.link.y {
					w.draft.n.model[w.cursor.x].stage[w.cursor.y].source = append(w.draft.n.model[w.cursor.x].stage[w.cursor.y].source[:i], w.draft.n.model[w.cursor.x].stage[w.cursor.y].source[i+1:]...)
					found = true
					break
				}
			}

			if !found && w.link.x < w.cursor.x && len(w.draft.n.model[w.link.x].stage) > w.link.y {
				w.draft.n.model[w.cursor.x].stage[w.cursor.y].source = append(w.draft.n.model[w.cursor.x].stage[w.cursor.y].source, projectModelSource{s: w.link.x, n: w.link.y})
			}

			w.draft.n.measure()
		}
	} else {
		if len(w.draft.n.model[w.cursor.x].stage) == 0 {
			return nil
		}
		w.link = position{w.cursor.x, w.cursor.y}
		w.mode = editLink
	}

	return nil
}

func (w *editWidget) editNode(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	w.mode = editNode

	return app.openModal(newNodeSizeBox())
}

func (w *editWidget) deleteNode(_ *gocui.Gui, _ *gocui.View) error {

	if w.draft == nil || !w.draft.ed {
		return nil
	}

	w.mode = editNode

	if w.cursor.x > 0 && w.cursor.x < len(w.draft.n.model)-1 {
		w.draft.n.model[w.cursor.x].removeNode(w.cursor.y)
		for i := range w.draft.n.model {
			for ii := range w.draft.n.model[i].stage {
				for iii, src := range w.draft.n.model[i].stage[ii].source {
					if src.s == w.cursor.x && src.n == w.cursor.y {
						w.draft.n.model[i].stage[ii].source = append(w.draft.n.model[i].stage[ii].source[:iii], w.draft.n.model[i].stage[ii].source[iii+1:]...)
						break
					}
				}
			}
		}
		w.shiftCursorY()
		w.draft.n.measure()
	}

	return nil
}

func (w *editWidget) linkEsc(_ *gocui.Gui, _ *gocui.View) error {

	w.mode = editNode

	return nil
}

func (w *editWidget) hideInvalid() {

	time.Sleep(1 * time.Second)
	w.setInvalid(false)
	gui.Update(func(g *gocui.Gui) error {
		return nil
	})

}

func (w *editWidget) headerTop() string {

	d := "  "
	if w.draft.ed {
		d = " \033[32;1m*\033[0m"
	}

	id := "## ?"
	if w.draft.status != psNew {
		id = fmt.Sprintf("## %v", w.draft.id)
	}

	o := ""
	if w.draft.o != nil {
		o = fmt.Sprintf(" \033[30;1m[%v]\033[0m", w.draft.o.id)
	}

	return fmt.Sprintf(" %s  %s  %s%s", d, w.draft.status.text(), id, o)
}

func (w *editWidget) headerBottom() string {

	sn := ""
	if len(w.draft.n.model[w.cursor.x].stage) > 0 {
		sn = fmt.Sprintf("  stage = %v : %v  node = %v : %v",
			w.draft.n.model[w.cursor.x].size, w.draft.n.model[w.cursor.x].volume,
			w.draft.n.model[w.cursor.x].stage[w.cursor.y].size, w.draft.n.model[w.cursor.x].stage[w.cursor.y].volume,
		)
	}

	return fmt.Sprintf("     total = %v : %v%s", w.draft.n.size, w.draft.n.volume, sn)
}

func (w *editWidget) body(buf []string, x, y int) []string {

	width := int(math.Floor((float64(x) - 10.0) / 5.0))
	if x-10-width*5 >= 2 {
		width++
	}
	height := int(math.Floor((float64(y) - 12.0) / 2.0))
	if y-12-height*2 >= 1 {
		height++
	}

	buf = w.bodyTop(buf, width, x)
	buf = w.bodyCenter(buf, width, height)
	buf = w.bodyBottom(buf, width, height, x)

	return buf
}

func (w *editWidget) bodyTop(buf []string, width, x int) []string {

	line := ""
	span := x - 10

	for i := range min(len(w.draft.n.model), width) {

		index := w.getX() + i

		color := "7;4"
		if !w.links() && i == w.cursor.x {
			color = "5;4"
		}
		if w.links() && i == w.link.x {
			color = "3;4"
		}
		if w.invalid && len(w.draft.n.model[index].stage) == 0 {
			color = "1;7"
		}

		size := 5
		if span < 5 {
			size = span
		}
		s := trail(index, size, 0)
		if index == 0 {
			s = trail("in", size, 0)
		}
		if index == len(w.draft.n.model)-1 {
			s = trail("out", size, 0)
		}
		line += fmt.Sprintf("\033[3%sm%s\033[0m", color, s)
		span -= size

	}

	buf = append(buf, fmt.Sprintf("  \033[30;1m┌\033[0m  %s  \033[30;1m┐\033[0m  ", space(line, 0, span)))

	scroll := " "
	if w.getY() > 0 {
		scroll = "▲"
	}

	buf = append(buf, fmt.Sprintf("     %s  \033[30;1m%s\033[0m  ", space("", x-10, 0), scroll))

	return buf
}

func (w *editWidget) bodyCenter(buf []string, width, height int) []string {

	grid := make([]string, height)

	for i := range min(len(w.draft.n.model), width) {
		s := w.draft.n.model[w.getX()+i]
		yy := 0

		for ii := w.getY(); ii < min(len(s.stage), height); ii++ {
			n := s.stage[ii]

			color := "7;4"
			for _, src := range w.draft.n.model[w.cursor.x+w.offset.x].stage[w.cursor.y+w.offset.y].source {
				if src.s == w.getX()+i && src.n == w.getY()+yy {
					color = "2;4"
				}
			}
			if i == w.cursor.x && yy == w.cursor.y {
				color = "5;4"
			}
			if w.invalid && !n.valid {
				color = "1;1"
			}

			char := "██"
			if w.links() && i == w.link.x && yy == w.link.y {
				char = "▓▓"
			}

			grid[yy] += fmt.Sprintf("\033[3%sm%s\033[0m   ", color, char)
			yy++
		}

		h := height - yy
		for range h {
			grid[yy] += space("", 0, 5)
			yy++
		}

	}

	for i, line := range grid {

		color := "7;4"
		if !w.links() && i == w.cursor.y {
			color = "5;4"
		}
		if w.links() && i == w.link.y {
			color = "3;4"
		}

		num := lead(w.getY()+i+1, 3, 0)
		if len(line) == min(len(w.draft.n.model), width)*5 {
			num = lead("", 3, 0)
		}

		buf = append(buf, fmt.Sprintf("\033[3%sm%s\033[0m  %s", color, num, line))

		if i != len(grid)-1 {
			buf = append(buf, " ")
		}
	}

	return buf
}

func (w *editWidget) bodyBottom(buf []string, width, height, x int) []string {

	s := 0
	for i := range min(len(w.draft.n.model), width) {
		s = max(len(w.draft.n.model[w.getX()+i].stage), s)
	}
	scroll := " "
	if (!w.links() && s-w.offset.y > height) || (w.links() && s-w.scroll.y > height) {
		scroll = "▼"
	}

	buf = append(buf, fmt.Sprintf("     %s  \033[30;1m%s\033[0m  ", space("", x-10, 0), scroll))

	left := " "
	if w.getX() > 0 {
		left = "◄"
	}

	right := " "
	if (!w.links() && len(w.draft.n.model)-w.offset.x > width) || (w.links() && len(w.draft.n.model)-w.scroll.x > width) {
		right = "►"
	}

	buf = append(buf, fmt.Sprintf("  \033[30;1m└%s %s %s┘\033[0m", left, space("", x-10, 0), right))

	return buf
}

func (w *editWidget) selectStageLeft() {

	if w.links() {
		if w.link.x > 0 {
			w.setLink(position{w.link.x - 1, w.link.y})
		}
	} else {
		if w.cursor.x > 0 {
			w.setCursor(position{w.cursor.x - 1, w.cursor.y})
		}
	}

}

func (w *editWidget) selectStageRight() {

	if w.links() {
		if w.link.x < len(w.draft.n.model)-1 {
			w.setLink(position{w.link.x + 1, w.link.y})
		}
	} else {
		if w.cursor.x < len(w.draft.n.model)-1 {
			w.setCursor(position{w.cursor.x + 1, w.cursor.y})
		}
	}

}

func (w *editWidget) selectNodeUp() {

	if w.links() {
		if w.link.y > 0 {
			w.link.y--
		}
	} else {
		if w.cursor.y > 0 {
			w.cursor.y--
		}
	}

}

func (w *editWidget) selectNodeDown() {

	if w.links() {
		if w.link.y < len(w.draft.n.model[w.link.x].stage)-1 {
			w.link.y++
		}
	} else {
		if w.cursor.y < len(w.draft.n.model[w.cursor.x].stage)-1 {
			w.cursor.y++
		}
	}

}

func (w *editWidget) shiftCursorY() {

	if w.links() {
		if w.link.y != 0 && w.link.y > len(w.draft.n.model[w.link.x].stage)-1 {
			w.link.y = max(0, len(w.draft.n.model[w.link.x].stage)-1)
		}
	} else {
		if w.cursor.y != 0 && w.cursor.y > len(w.draft.n.model[w.cursor.x].stage)-1 {
			w.cursor.y = max(0, len(w.draft.n.model[w.cursor.x].stage)-1)
		}
	}

}

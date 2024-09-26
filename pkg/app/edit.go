// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"math"
	"slices"
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
		w.link = v
		w.mark()
	}

}

func (w *editWidget) setScroll(v position) {

	if v.x != w.scroll.x || v.y != w.scroll.y {
		w.scroll = v
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

func (w *editWidget) edits() bool {

	if w.draft == nil || w.draft.status != projectEdit {
		return false
	}

	return true
}

func (w *editWidget) setNodeSize(v int) {

	if v != w.draft.getNodeSize(w.cursor.x, w.cursor.y) {
		w.draft.setNodeSize(w.cursor.x, w.cursor.y, v)
		w.mark()
		w.draft.measure()
	}

}

func (w *editWidget) getNodeSize() int {

	return w.draft.getNodeSize(w.cursor.x, w.cursor.y)
}

func newEditWidget() *editWidget {

	w := &editWidget{
		draft:   nil,
		mode:    editNode,
		invalid: false,
		link:    newPosition(0, 0),
		scroll:  newPosition(0, 0),
	}
	w.widget = newWidget(w, "edit")

	return w
}

func (w *editWidget) transform(x int, y int) (int, int, int, int) {

	return -1, 0, x, y - 1
}

func (w *editWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	buf := make([]string, 5)

	if w.draft != nil {

		x, y := gui.Size()

		width := int(math.Floor((float64(x) - 10.0) / 5.0))
		if x-10-width*5 >= 3 {
			width++
		}
		height := int(math.Floor((float64(y) - 12.0) / 2.0))
		if y-12-height*2 == 1 {
			height++
		}

		if w.links() {
			w.setScroll(w.move(w.link, w.scroll, width, height))
		} else {
			w.setOffset(w.move(w.cursor, w.offset, width, height))
		}

		buf[0] = " "
		buf[1] = w.headerTop()
		buf[2] = " "
		buf[3] = w.headerBottom()
		buf[4] = " "

		buf = w.body(buf, x, width, height)

	}

	return buf, nil
}

func (w *editWidget) keybinding() error {

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlN, gocui.ModNone, w.newProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlE, gocui.ModNone, w.editProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlD, gocui.ModNone, w.dubProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlS, gocui.ModNone, w.saveProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlV, gocui.ModNone, w.validateProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.keyArrowLeft); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.keyArrowRight); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.keyArrowUp); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.keyArrowDown); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyInsert, gocui.ModNone, w.insertStage); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyDelete, gocui.ModNone, w.deleteStage); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.insertNode); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.editNodeSource); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlSpace, gocui.ModNone, w.editNode); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyBackspace, gocui.ModNone, w.deleteNode); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.cancel); err != nil {
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
		w.mark()
	}

	return nil
}

func (w *editWidget) dubProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	p := app.v.result.getFocused()
	if p != nil {
		w.setDraft(newProject(p))
		w.reset()
		w.mark()
	}

	return nil
}

func (w *editWidget) saveProject(_ *gocui.Gui, _ *gocui.View) error {

	w.setMode(editNode)

	if w.draft != nil {
		w.draft.save()
		w.mark()
		app.v.result.mark()
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

	if !w.edits() {
		return nil
	}

	w.setMode(editNode)
	w.selectStageRight()
	w.draft.insertStage(w.cursor.x)
	w.shiftCursorY()
	w.mark()

	return nil
}

func (w *editWidget) deleteStage(_ *gocui.Gui, _ *gocui.View) error {

	if !w.edits() {
		return nil
	}

	w.setMode(editNode)
	w.draft.deleteStage(w.cursor.x)
	w.shiftCursorY()
	w.draft.measure()
	w.mark()

	return nil
}

func (w *editWidget) insertNode(_ *gocui.Gui, _ *gocui.View) error {

	if !w.edits() {
		return nil
	}

	w.setMode(editNode)

	if w.cursor.x > 0 && w.cursor.x < len(w.draft.getModel())-1 {
		w.draft.getStage(w.cursor.x).addModel(newNode(0, nil))
		w.setCursor(newPosition(w.cursor.x, len(w.draft.getStage(w.cursor.x).model)-1))
		w.mark()
		return app.openModal(newNodeSizeBox())
	}

	return nil
}

func (w *editWidget) editNodeSource(_ *gocui.Gui, _ *gocui.View) error {

	if !w.edits() {
		return nil
	}

	if w.links() {

		if w.cursor.x != w.link.x || w.cursor.y != w.link.y {

			found := false
			for i, src := range w.draft.getNode(w.cursor).source {
				if src.x == w.link.x && src.y == w.link.y {
					w.draft.getNode(w.cursor).source = slices.Delete(w.draft.getNode(w.cursor).source, i, i+1)
					found = true
					break
				}
			}

			if !found && w.link.x < w.cursor.x && len(w.draft.getStage(w.link.x).model) > w.link.y {
				w.draft.getNode(w.cursor).source = append(w.draft.getNode(w.cursor).source, w.link)
			}

			w.draft.measure()
		}
		w.setMode(editNode)

	} else {

		if len(w.draft.getStage(w.cursor.x).model) == 0 || w.cursor.x == 0 {
			return nil
		}
		w.setLink(w.cursor)
		w.setScroll(w.offset)
		w.setMode(editLink)

	}
	w.mark()

	return nil
}

func (w *editWidget) editNode(_ *gocui.Gui, _ *gocui.View) error {

	if !w.edits() {
		return nil
	}

	w.setMode(editNode)

	return app.openModal(newNodeSizeBox())
}

func (w *editWidget) deleteNode(_ *gocui.Gui, _ *gocui.View) error {

	if !w.edits() {
		return nil
	}

	w.setMode(editNode)

	if w.cursor.x > 0 && w.cursor.x < len(w.draft.getModel())-1 {
		w.draft.getStage(w.cursor.x).deleteModel(w.cursor.y)
		for _, s := range w.draft.getModel() {
			for _, n := range s.model {
				for iii, src := range n.source {
					if src.x == w.cursor.x && src.y == w.cursor.y {
						n.source = slices.Delete(n.source, iii, iii+1)
						break
					}
				}
			}
		}
		w.shiftCursorY()
		w.draft.measure()
		w.mark()
	}

	return nil
}

func (w *editWidget) cancel(_ *gocui.Gui, _ *gocui.View) error {

	if w.links() {

		w.setMode(editNode)

	} else {

		if !w.edits() {
			return nil
		}

		if w.draft != nil {
			if w.draft.model != nil {
				w.draft.cancel()
				w.mark()
				app.v.result.mark()
				return nil
			}
			w.draft = nil
			w.mark()
		}

	}

	return nil
}

func (w *editWidget) hideInvalid() {

	time.Sleep(1 * time.Second)
	w.setInvalid(false)
	app.update(nil)

}

func (w *editWidget) headerTop() string {

	d := "  "
	if w.edits() {
		d = " \033[32;1m*\033[0m"
	}

	return fmt.Sprintf(" %s  %s  ## %v", d, w.draft.status.text(), w.draft.id)
}

func (w *editWidget) headerBottom() string {

	sn := ""
	if len(w.draft.getStage(w.cursor.x).model) > 0 {
		sn = fmt.Sprintf("  stage = %v : %v  node = %v : %v",
			w.draft.getStage(w.cursor.x).size, w.draft.getStage(w.cursor.x).volume,
			w.draft.getNode(w.cursor).size, w.draft.getNode(w.cursor).volume,
		)
	}

	return fmt.Sprintf("     total = %v : %v%s", w.draft.getSize(), w.draft.getVolume(), sn)
}

func (w *editWidget) body(buf []string, x, width, height int) []string {

	buf = w.bodyTop(buf, width, x)
	buf = w.bodyCenter(buf, width, height)
	buf = w.bodyBottom(buf, width, height, x)

	return buf
}

func (w *editWidget) bodyTop(buf []string, width, x int) []string {

	line := ""
	span := x - 10

	for i := range min(len(w.draft.getModel()), width) {

		index := w.getX() + i

		color := "7;4"
		if !w.links() && index == w.cursor.x {
			color = "5;4"
		}
		if w.links() && index == w.link.x {
			color = "3;4"
		}
		if w.invalid && len(w.draft.getStage(index).model) == 0 {
			color = "1;7"
		}

		size := 5
		if span < 5 {
			size = span
		}
		str := trail(index, size, 0)
		if index == 0 {
			str = trail("in", size, 0)
		}
		if index == len(w.draft.getModel())-1 {
			str = trail("out", size, 0)
		}
		line += fmt.Sprintf("\033[3%sm%s\033[0m", color, str)
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
	cursrc := []position{}

	if len(w.draft.getStage(w.cursor.x).model) > 0 {
		cursrc = w.draft.getNode(newPosition(w.cursor.x, w.cursor.y)).source
	}

	for i := range min(len(w.draft.getModel()), width) {
		s := w.draft.getStage(w.getX() + i)
		yy := 0

		for ii := w.getY(); ii < min(len(s.model), height+w.getY()); ii++ {
			n := s.model[ii]

			color := "7;4"
			for _, src := range cursrc {
				if src.x == w.getX()+i && src.y == ii {
					color = "2;4"
				}
			}
			if w.getX()+i == w.cursor.x && ii == w.cursor.y {
				color = "5;4"
			}
			if w.invalid && !n.valid {
				color = "1;1"
			}

			char := "██"
			if w.links() && w.getX()+i == w.link.x && ii == w.link.y {
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
		if !w.links() && i+w.offset.y == w.cursor.y {
			color = "5;4"
		}
		if w.links() && i+w.scroll.y == w.link.y {
			color = "3;4"
		}

		num := lead(w.getY()+i+1, 3, 0)
		if len(line) == min(len(w.draft.getModel()), width)*5 {
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
	for i := range min(len(w.draft.getModel()), width) {
		s = max(len(w.draft.getStage(w.getX()+i).model), s)
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
	if (!w.links() && len(w.draft.getModel())-w.offset.x > width) || (w.links() && len(w.draft.getModel())-w.scroll.x > width) {
		right = "►"
	}

	buf = append(buf, fmt.Sprintf("  \033[30;1m└%s %s %s┘\033[0m", left, space("", x-10, 0), right))

	return buf
}

func (w *editWidget) selectStageLeft() {

	if w.links() {
		if w.link.x > 0 {
			w.setLink(newPosition(w.link.x-1, w.link.y))
		}
	} else {
		if w.cursor.x > 0 {
			w.setCursor(newPosition(w.cursor.x-1, w.cursor.y))
		}
	}

}

func (w *editWidget) selectStageRight() {

	if w.links() {
		if w.link.x < len(w.draft.getModel())-1 {
			w.setLink(newPosition(w.link.x+1, w.link.y))
		}
	} else {
		if w.cursor.x < len(w.draft.getModel())-1 {
			w.setCursor(newPosition(w.cursor.x+1, w.cursor.y))
		}
	}

}

func (w *editWidget) selectNodeUp() {

	if w.links() {
		if w.link.y > 0 {
			w.setLink(newPosition(w.link.x, w.link.y-1))
		}
	} else {
		if w.cursor.y > 0 {
			w.setCursor(newPosition(w.cursor.x, w.cursor.y-1))
		}
	}

}

func (w *editWidget) selectNodeDown() {

	if w.links() {
		if w.link.y < len(w.draft.getStage(w.link.x).model)-1 {
			w.setLink(newPosition(w.link.x, w.link.y+1))
		}
	} else {
		if w.cursor.y < len(w.draft.getStage(w.cursor.x).model)-1 {
			w.setCursor(newPosition(w.cursor.x, w.cursor.y+1))
		}
	}

}

func (w *editWidget) shiftCursorY() {

	if w.links() {
		if w.link.y != 0 && w.link.y > len(w.draft.getStage(w.link.x).model)-1 {
			w.setLink(newPosition(w.link.x, max(0, len(w.draft.getStage(w.link.x).model)-1)))
		}
	} else {
		if w.cursor.y != 0 && w.cursor.y > len(w.draft.getStage(w.cursor.x).model)-1 {
			w.setCursor(newPosition(w.cursor.x, max(0, len(w.draft.getStage(w.cursor.x).model)-1)))
		}
	}

}

func (w *editWidget) move(cursor, offset position, width, height int) position {

	if cursor.x >= offset.x+width {
		offset.x = cursor.x - width + 1
	}
	if cursor.x < offset.x {
		offset.x = cursor.x
	}
	if cursor.y >= offset.y+height {
		offset.y = cursor.y - height + 1
	}
	if cursor.y < offset.y {
		offset.y = cursor.y
	}

	return offset
}

package app

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type resultWidget struct {
	*widget
	holded bool
	show   []*project
	width  map[string]int
}

func (w *resultWidget) setHolded(v bool) {

	if v != w.holded {
		w.holded = v
		w.mark()
	}

}

func (w *resultWidget) getActivated() int {

	cnt := 0
	for _, p := range w.show {
		if p.status == projectActive {
			cnt++
		}
	}

	return cnt
}

func (w *resultWidget) getFocused() *project {

	if w.cursor.y < len(w.show) {
		return w.show[w.cursor.y]
	}

	return nil
}

func (w *resultWidget) getSelected() bool {

	for _, p := range w.show {
		if p.selected {
			return true
		}
	}

	return false
}

func (w *resultWidget) updateShow(height int) {

	l := make([]*project, 0, len(app.s.ev))
	for _, p := range app.s.ev {
		if ((!w.holded && p.status == projectActive) || (w.holded && p.status != projectTerminated)) && p.goal {
			l = append(l, p)
		}
	}
	for _, p := range app.s.ev {
		if ((!w.holded && p.status == projectActive) || (w.holded && p.status != projectTerminated)) && !p.goal {
			l = append(l, p)
		}
	}
	w.show = l

	w.shift(height)

}

func (w *resultWidget) updateWidth() {

	for key := range w.width {
		w.width[key] = 0
	}

	for _, p := range w.show {
		w.width["num"] = max(len(fmt.Sprint(p.id)), w.width["num"])
		w.width["status"] = max(len(p.status.text()), w.width["status"])
		w.width["size"] = max(len(fmt.Sprint(p.stat.size)), w.width["size"])
		w.width["volume"] = max(len(fmt.Sprint(p.stat.volume)), w.width["volume"])
		w.width["gen"] = max(len(fmt.Sprint(p.stat.gen)), w.width["gen"])
		w.width["ev"] = max(len(fmt.Sprint(p.stat.ev)), w.width["ev"])
		w.width["age"] = max(len(fmt.Sprint(p.stat.age)), w.width["age"])
		w.width["best"] = max(len(p.stat.best), w.width["best"])
	}

	w.width["num"] = max(w.width["num"]+2, 7)
	w.width["status"] = max(w.width["status"]+2, 10)
	w.width["size"] = max(w.width["size"]+2, 9)
	w.width["volume"] = max(w.width["volume"]+2, 11)
	w.width["gen"] = max(w.width["gen"]+2, 8)
	w.width["ev"] = max(w.width["ev"]+2, 8)
	w.width["age"] = max(w.width["age"]+2, 8)
	w.width["best"] = max(w.width["best"]+2, 9)

}

func newResultWidget() *resultWidget {

	w := &resultWidget{
		holded: false,
		show:   make([]*project, 0),
		width: map[string]int{
			"num":    0,
			"status": 0,
			"size":   0,
			"volume": 0,
			"gen":    0,
			"ev":     0,
			"age":    0,
			"best":   0,
		},
	}
	w.widget = newWidget(w, "result")

	return w
}

func (w *resultWidget) transform(x int, y int) (int, int, int, int) {

	return -1, 0, x, y - 1
}

func (w *resultWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, y := gui.Size()

	height := int(math.Floor((float64(y) - 6.0) / 2.0))

	w.updateShow(height)
	w.updateWidth()

	if w.cursor.y > w.offset.y+height-1 {
		w.setOffset(newPosition(0, w.cursor.y-height+1))
	}
	if w.cursor.y < w.offset.y {
		w.setOffset(newPosition(0, w.cursor.y))
	}

	buf := make([]string, y-2)

	buf[0] = " "
	buf[1] = w.header()
	buf = w.body(buf, height, x, y)

	return buf, nil
}

func (w *resultWidget) keybinding() error {

	if err = gui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.runProd); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.stopProd); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.moveDown); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.moveUp); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlN, gocui.ModNone, w.newProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlE, gocui.ModNone, w.editProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlD, gocui.ModNone, w.dubProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.selectProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlSpace, gocui.ModNone, w.selectAllProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlH, gocui.ModNone, w.holdProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlA, gocui.ModNone, w.activateProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlT, gocui.ModNone, w.terminateProject); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyTab, gocui.ModNone, w.toggleHolded); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlS, gocui.ModNone, w.saveLab); err != nil {
		return err
	}

	if err = gui.SetKeybinding(w.name, gocui.KeyCtrlL, gocui.ModNone, w.loadLab); err != nil {
		return err
	}

	return nil
}

func (w *resultWidget) runProd(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.getActivated() == 0 {
			return nil
		}
		app.apply(appRun, "Run...")
		w.mark()
		app.s.lab.Run()
		go w.checkStats()

	}

	return nil
}

func (w *resultWidget) stopProd(_ *gocui.Gui, _ *gocui.View) error {

	if app.s.status == appRun {

		app.apply(appWait, "Wait for stop examine")
		w.mark()
		app.s.lab.Stop()
		go w.checkStatus()

	}

	return nil
}

func (w *resultWidget) moveDown(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.cursor.y < len(w.show)-1 {
			w.setCursor(newPosition(0, w.cursor.y+1))
		}

	} else {

		_, y := gui.Size()
		yy := w.offset.y + int(math.Floor(float64(y)-4.0)/2.0)

		if yy >= len(w.show)-1 {
			w.setCursor(newPosition(0, len(w.show)-1))
		} else {
			w.setCursor(newPosition(0, yy))
			w.setOffset(newPosition(0, w.offset.y+1))
		}

	}

	return nil
}

func (w *resultWidget) moveUp(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.cursor.y > 0 {
			w.setCursor(newPosition(0, w.cursor.y-1))
		}

	} else {

		if w.offset.y > 0 {
			w.setCursor(newPosition(0, w.offset.y-1))
			w.setOffset(newPosition(0, w.offset.y))
		}

	}

	return nil
}

func (w *resultWidget) newProject(_ *gocui.Gui, _ *gocui.View) error {

	app.v.edit.setDraft(newProject(nil))
	app.v.edit.reset()

	return app.showEdit(nil, nil)
}

func (w *resultWidget) editProject(_ *gocui.Gui, _ *gocui.View) error {

	p := w.getFocused()
	if p != nil {
		app.v.edit.setDraft(p)
		app.v.edit.draft.edit()
		app.v.edit.reset()
		w.mark()

		return app.showEdit(nil, nil)
	}

	return nil
}

func (w *resultWidget) dubProject(_ *gocui.Gui, _ *gocui.View) error {

	p := w.getFocused()
	if p != nil {
		app.v.edit.setDraft(newProject(p))
		app.v.edit.reset()

		return app.showEdit(nil, nil)
	}

	return nil
}

func (w *resultWidget) selectProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	p := w.getFocused()
	if p != nil {
		p.selected = !p.selected
		w.mark()
	}

	return nil
}

func (w *resultWidget) selectAllProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	s := !w.getSelected()
	for _, p := range w.show {
		p.selected = s
		w.mark()
	}

	return nil
}

func (w *resultWidget) holdProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if w.getSelected() {
		for _, p := range w.show {
			if p.status == projectActive && p.selected {
				p.setStatus(projectHolded)
				w.mark()
			}
		}
	} else {
		p := w.getFocused()
		if p != nil && p.status == projectActive {
			p.setStatus(projectHolded)
			w.mark()
		}
	}

	if !w.holded && w.getSelected() {
		for _, p := range w.show {
			p.selected = false
		}
		w.mark()
	}

	return nil
}

func (w *resultWidget) activateProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if !w.holded {
		return nil
	}

	if w.getSelected() {
		for _, p := range w.show {
			if p.status == projectHolded && p.selected {
				p.setStatus(projectActive)
				w.mark()
			}
		}
	} else {
		p := w.getFocused()
		if p != nil && p.status == projectHolded {
			p.setStatus(projectActive)
			w.mark()
		}
	}

	return nil
}

func (w *resultWidget) terminateProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if w.getSelected() {
		for _, p := range w.show {
			if p.selected {
				p.selected = false
				p.setStatus(projectTerminated)
				w.mark()
			}
		}
	} else {
		p := w.getFocused()
		if p != nil {
			p.selected = false
			p.setStatus(projectTerminated)
			w.mark()
		}
	}

	if !w.holded && w.getSelected() {
		for _, p := range w.show {
			p.selected = false
		}
		w.mark()
	}

	return nil
}

func (w *resultWidget) toggleHolded(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	w.setHolded(!w.holded)

	return nil
}

func (w *resultWidget) saveLab(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	app.apply(appWait, "Saving lab data...")
	go app.saveLab()

	return nil
}

func (w *resultWidget) loadLab(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	app.apply(appWait, "Loading lab data...")
	go app.loadLab()

	return nil
}

func (w *resultWidget) checkStatus() {

	time.Sleep(1 * time.Second)

	if !app.s.lab.GetExec() {
		app.apply(appIdle, "")
		w.mark()
		app.update(nil)
	} else {
		go w.checkStatus()
	}

}

func (w *resultWidget) checkStats() {

	time.Sleep(1 * time.Second)

	if app.s.lab.GetExec() {
		go w.checkStats()
	} else {
		w.stopProd(nil, nil)
	}

	for _, p := range app.s.ev {
		if p.status != projectTerminated {
			p.refine()
		}
	}

	w.mark()
	app.update(nil)

}

func (w *resultWidget) header() string {

	t := fmt.Sprintf("  %s│", space("##", 1, w.width["num"]-3))
	if w.holded {
		t += fmt.Sprintf("%s│", space("\033[33;4mSTATUS\033[0m", 1, w.width["status"]-7))
	} else {
		t += fmt.Sprintf("%s│", space("STATUS", 1, w.width["status"]-7))
	}
	t += fmt.Sprintf("%s│", space("SIZE", w.width["size"]-5, 1))
	t += fmt.Sprintf("%s│", space("VOLUME", w.width["volume"]-7, 1))
	t += fmt.Sprintf("%s│", space("GEN", w.width["gen"]-4, 1))
	t += fmt.Sprintf("%s│", space("EV", w.width["ev"]-3, 1))
	t += fmt.Sprintf("%s│", space("AGE", w.width["age"]-4, 1))
	t += fmt.Sprintf("%s", space("BEST", w.width["best"]-5, 1))

	return t
}

func (w *resultWidget) body(buf []string, height, x, y int) []string {

	for i := range min(len(w.show), height) {

		buf[2+i*2+0] = fmt.Sprintf("  %s┼%s┼%s┼%s┼%s┼%s┼%s┼%s", strings.Repeat("─", w.width["num"]), strings.Repeat("─", w.width["status"]), strings.Repeat("─", w.width["size"]), strings.Repeat("─", w.width["volume"]), strings.Repeat("─", w.width["gen"]), strings.Repeat("─", w.width["ev"]), strings.Repeat("─", w.width["age"]), strings.Repeat("─", w.width["best"]))

		p := w.show[w.offset.y+i]
		textColor, statusColor := w.color(p, w.cursor.y == w.offset.y+i)
		buf[2+i*2+1] = fmt.Sprintf("%s \033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m", w.selectedMark(p), textColor, trail(p.id, w.width["num"], 1), statusColor, trail(p.status.text(), w.width["status"], 1), textColor, lead(p.stat.size, w.width["size"], 1), textColor, lead(p.stat.volume, w.width["volume"], 1), textColor, lead(p.stat.gen, w.width["gen"], 1), textColor, lead(p.stat.ev, w.width["ev"], 1), textColor, lead(p.stat.age, w.width["age"], 1), textColor, lead(p.stat.best, w.width["best"], 1))

	}

	buf[2+min(len(w.show), height)*2] = fmt.Sprintf("  %s┴%s┴%s┴%s┴%s┴%s┴%s┴%s", strings.Repeat("─", w.width["num"]), strings.Repeat("─", w.width["status"]), strings.Repeat("─", w.width["size"]), strings.Repeat("─", w.width["volume"]), strings.Repeat("─", w.width["gen"]), strings.Repeat("─", w.width["ev"]), strings.Repeat("─", w.width["age"]), strings.Repeat("─", w.width["best"]))

	if w.offset.y > 0 {
		buf[0] = space("▲", x-1, 0)
	}

	if len(w.show)-w.offset.y > height {
		buf[y-3] = space("▼", x-1, 0)
	}

	return buf
}

func (w *resultWidget) selectedMark(p *project) string {

	if p.selected {
		return "█"
	}

	return " "
}

func (w *resultWidget) color(p *project, focused bool) (string, string) {

	if focused && app.idle() {
		if p.stat.goal {
			return "2;7", "2;7"
		}
		return "7;7", "7;7"
	}

	text := "0;1"
	status := "1;1"
	if p.status == projectActive {
		text = "7;4"
		if p.stat.goal {
			status = "2;1"
		} else {
			status = "7;4"
		}
	} else if p.status == projectHolded {
		text = "0;1"
		status = "3;4"
	}
	if p.stat.goal {
		text = "2;1"
	}

	return text, status
}

func (w *resultWidget) shift(height int) {

	if w.cursor.y > 0 && w.cursor.y > len(w.show)-1 {

		w.setCursor(newPosition(0, len(w.show)-1))
		w.setOffset(newPosition(0, len(w.show)-min(len(w.show), height)))

	}

}

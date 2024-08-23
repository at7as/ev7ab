package app

import (
	"fmt"
	"math"
	"strings"

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

func (w *resultWidget) updateShow() {

	l := make([]*project, 0, len(app.s.list))
	for _, p := range app.s.list {
		if (!w.holded && p.status == psActive) || (w.holded && p.status != psTerminated) {
			l = append(l, p)
		}
	}
	w.show = l

	w.mark()

}

func (w *resultWidget) updateWidth() {

	for key := range w.width {
		w.width[key] = 0
	}

	for _, p := range w.show {
		w.width["num"] = max(len(fmt.Sprint(p.id)), w.width["num"])
		w.width["status"] = max(len(p.status.text()), w.width["status"])
		w.width["size"] = max(len(fmt.Sprint(p.s.size)), w.width["size"])
		w.width["volume"] = max(len(fmt.Sprint(p.s.volume)), w.width["volume"])
		w.width["gen"] = max(len(fmt.Sprint(p.s.gen)), w.width["gen"])
		w.width["ev"] = max(len(fmt.Sprint(p.s.ev)), w.width["ev"])
		w.width["age"] = max(len(fmt.Sprint(p.s.age)), w.width["age"])
		w.width["best"] = max(len(p.s.best), w.width["best"])
	}

	w.width["num"] = max(w.width["num"]+2, 7)
	w.width["status"] = max(w.width["status"]+2, 10)
	w.width["size"] = max(w.width["size"]+2, 9)
	w.width["volume"] = max(w.width["volume"]+2, 11)
	w.width["gen"] = max(w.width["gen"]+2, 8)
	w.width["ev"] = max(w.width["ev"]+2, 8)
	w.width["age"] = max(w.width["age"]+2, 8)
	w.width["best"] = max(w.width["best"]+2, 9)

	w.mark()

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

// top to result when project is goal
func (w *resultWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, y := gui.Size()

	height := int(math.Floor((float64(y) - 6.0) / 2.0))

	if w.cursor.y > w.offset.y+height-1 {
		w.setOffset(position{0, w.cursor.y - height + 1})
	}
	if w.cursor.y < w.offset.y {
		w.setOffset(position{0, w.cursor.y})
	}

	w.updateWidth()
	w.updateShow()

	buf := make([]string, y-2)

	buf[0] = " "
	buf[1] = w.header()
	buf = w.body(buf, height, x, y)

	return buf, nil
}

func (w *resultWidget) keybinding() error {

	if err := gui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.toggleRunApp); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.moveDown); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.moveUp); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlN, gocui.ModNone, w.newProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlE, gocui.ModNone, w.editProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlD, gocui.ModNone, w.dubProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.selectProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlSpace, gocui.ModNone, w.selectAllProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlH, gocui.ModNone, w.holdProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlA, gocui.ModNone, w.activateProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyCtrlT, gocui.ModNone, w.terminateProject); err != nil {
		return err
	}

	if err := gui.SetKeybinding(w.name, gocui.KeyTab, gocui.ModNone, w.toggleHolded); err != nil {
		return err
	}

	return nil
}

func (w *resultWidget) toggleRunApp(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.activated() == 0 {
			return nil
		}
		app.apply(appRun, "Run...")

	} else {

		app.apply(appIdle, "")

	}

	return nil
}

func (w *resultWidget) moveDown(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.cursor.y < len(w.show)-1 {
			w.setCursor(position{0, w.cursor.y + 1})
		}

	} else {

		_, y := gui.Size()
		yy := w.offset.y + int(math.Floor(float64(y)-4.0)/2.0)

		if yy >= len(w.show)-1 {
			w.setCursor(position{0, len(w.show) - 1})
		} else {
			w.setCursor(position{0, yy})
			w.setOffset(position{0, w.offset.y + 1})
		}

	}

	return nil
}

func (w *resultWidget) moveUp(_ *gocui.Gui, _ *gocui.View) error {

	if app.idle() {

		if w.cursor.y > 0 {
			w.setCursor(position{0, w.cursor.y - 1})
		}

	} else {

		if w.offset.y > 0 {
			w.setCursor(position{0, w.offset.y - 1})
			w.setOffset(position{0, w.offset.y})
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

	p := w.focused()
	if p != nil {
		app.v.edit.setDraft(p)
		app.v.edit.draft.edit()
		app.v.edit.reset()

		return app.showEdit(nil, nil)
	}

	return nil
}

func (w *resultWidget) dubProject(_ *gocui.Gui, _ *gocui.View) error {

	p := w.focused()
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

	p := w.focused()
	if p != nil {
		p.selected = !p.selected
	}

	return nil
}

func (w *resultWidget) selectAllProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	s := !w.selected()
	for _, p := range w.show {
		p.selected = s
	}

	return nil
}

func (w *resultWidget) holdProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if w.selected() {
		for _, p := range w.show {
			if p.status == psActive && p.selected {
				p.status = psHolded
			}
		}
	} else {
		p := w.focused()
		if p != nil && p.status == psActive {
			p.status = psHolded
		}
	}

	w.next()

	return nil
}

func (w *resultWidget) activateProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if !w.holded {
		return nil
	}

	if w.selected() {
		for _, p := range w.show {
			if p.status == psHolded && p.selected {
				p.status = psActive
			}
		}
	} else {
		p := w.focused()
		if p != nil && p.status == psHolded {
			p.status = psActive
		}
	}

	return nil
}

func (w *resultWidget) terminateProject(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	if w.selected() {
		for _, p := range w.show {
			if p.selected {
				p.selected = false
				p.status = psTerminated
			}
		}
	} else {
		p := w.focused()
		if p != nil {
			p.selected = false
			p.status = psTerminated
		}
	}

	w.next()

	return nil
}

func (w *resultWidget) toggleHolded(_ *gocui.Gui, _ *gocui.View) error {

	if !app.idle() {
		return nil
	}

	w.setHolded(!w.holded)

	w.next()

	return nil
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
		buf[2+i*2+1] = fmt.Sprintf("%s \033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m│\033[3%sm%s\033[0m", w.selectedMark(p), textColor, trail(p.id, w.width["num"], 1), statusColor, trail(p.status.text(), w.width["status"], 1), textColor, lead(p.s.size, w.width["size"], 1), textColor, lead(p.s.volume, w.width["volume"], 1), textColor, lead(p.s.gen, w.width["gen"], 1), textColor, lead(p.s.ev, w.width["ev"], 1), textColor, lead(p.s.age, w.width["age"], 1), textColor, lead(p.s.best, w.width["best"], 1))

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
		if p.s.goal {
			return "2;7", "2;7"
		}
		return "7;7", "7;7"
	}

	text := "0;1"
	status := "1;1"
	if p.status == psActive {
		text = "7;4"
		if p.s.goal {
			status = "2;1"
		} else {
			status = "7;4"
		}
	} else if p.status == psHolded {
		text = "0;1"
		status = "3;4"
	}
	if p.s.goal {
		text = "2;1"
	}

	return text, status
}

func (w *resultWidget) activated() int {

	cnt := 0
	for _, p := range w.show {
		if p.status == psActive {
			cnt++
		}
	}

	return cnt
}

func (w *resultWidget) focused() *project {

	if w.cursor.y < len(w.show) {
		return w.show[w.cursor.y]
	}

	return nil
}

func (w *resultWidget) selected() bool {

	for _, p := range w.show {
		if p.selected {
			return true
		}
	}

	return false
}

func (w *resultWidget) next() {

	y := w.cursor.y
	w.moveDown(nil, nil)
	if y == w.cursor.y {
		w.moveUp(nil, nil)
	}

}

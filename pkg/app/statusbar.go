package app

import (
	"fmt"
)

type statusbarWidget struct {
	*widget
	color int
	text  string
}

func (w *statusbarWidget) setColor(v int) {

	if v != w.color {
		w.color = v
		w.dig()
	}

}

func (w *statusbarWidget) getColor() string {

	switch w.color {
	case 0:
		return "7;7"
	case 1:
		return "5;7"
	case 2:
		return "3;7"
	}

	return "7;7"
}

func (w *statusbarWidget) setText(v string) {

	if v != w.text {
		w.text = v
		w.dig()
	}

}

func (w *statusbarWidget) getText() string {

	return w.text
}

func newStatusbarWidget() *statusbarWidget {

	w := &statusbarWidget{}
	w.widget = newWidget(w, "statusbar")

	return w
}

func (w *statusbarWidget) keybinding() error {

	return nil
}

func (w *statusbarWidget) transform(x int, y int) (int, int, int, int) {

	return -1, y - 2, x, y
}

func (w *statusbarWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	buf := []string{}

	x, _ := gui.Size()

	buf = append(buf,
		fmt.Sprintf("\033[3%sm%s\033[0m", w.getColor(), space(w.getText(), 1, x-1-len(w.getText()))),
	)

	return buf, nil
}

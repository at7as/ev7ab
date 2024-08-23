package app

import (
	"fmt"
)

type statusbarWidget struct {
	*widget
	color string
	text  string
}

func (w *statusbarWidget) setColor(s applicationStatus) {

	color := ""
	switch s {
	case appIdle:
		color = "7;7"
	case appRun:
		color = "5;7"
	case appWait:
		color = "3;7"
	}
	if color != w.color {
		w.color = color
		w.mark()
	}

}

func (w *statusbarWidget) setText(v string) {

	if v != w.text {
		w.text = v
		w.mark()
	}

}

func newStatusbarWidget() *statusbarWidget {

	w := &statusbarWidget{
		text: "",
	}
	w.widget = newWidget(w, "statusbar")
	w.setColor(appIdle)

	return w
}

func (w *statusbarWidget) transform(x int, y int) (int, int, int, int) {

	return -1, y - 2, x, y
}

func (w *statusbarWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, _ := gui.Size()

	buf := []string{}

	buf = append(buf,
		fmt.Sprintf("\033[3%sm%s\033[0m", w.color, space(w.text, 1, x-1-len(w.text))),
	)

	return buf, nil
}

func (w *statusbarWidget) keybinding() error {

	return nil
}

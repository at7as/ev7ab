package app

import (
	"fmt"
)

type keybarWidget struct {
	*widget
	tab keybarTab
}

func (w *keybarWidget) setTab(v keybarTab) {

	if v != w.tab {
		w.tab = v
		w.dig()
	}

}

func (w *keybarWidget) getTab() keybarTab {

	return w.tab
}

func newKeybarWidget() *keybarWidget {

	w := &keybarWidget{}
	w.widget = newWidget(w, "keybar")

	return w
}

func (w *keybarWidget) keybinding() error {

	return nil
}

func (w *keybarWidget) transform(x int, y int) (int, int, int, int) {

	return -1, -1, x, 1
}

func (w *keybarWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	buf := []string{}

	x, _ := gui.Size()

	buf = append(buf,
		fmt.Sprintf("\033[37;7m \033[3%sm Setup \033[37;7m \033[3%sm Result \033[37;7m \033[3%sm Edit \033[37;7m%s", w.color(tabSetup), w.color(tabResult), w.color(tabEdit), space("", 0, x-24)),
	)

	return buf, nil
}

func (w *keybarWidget) color(t keybarTab) string {

	if w.getTab() == t {
		return "2;7"
	}

	return "7;7"
}

type keybarTab string

const (
	tabSetup  keybarTab = "setup"
	tabResult keybarTab = "result"
	tabEdit   keybarTab = "edit"
)

// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
)

type keybarWidget struct {
	*widget
	tab keybarTab
}

type keybarTab string

const (
	tabSetup  keybarTab = "setup"
	tabResult keybarTab = "result"
	tabEdit   keybarTab = "edit"
)

func (w *keybarWidget) setTab(v keybarTab) {

	if v != w.tab {
		w.tab = v
		w.mark()
	}

}

func newKeybarWidget() *keybarWidget {

	w := &keybarWidget{
		tab: tabSetup,
	}
	w.widget = newWidget(w, "keybar")

	return w
}

func (w *keybarWidget) transform(x int, y int) (int, int, int, int) {

	return -1, -1, x, 1
}

func (w *keybarWidget) render() ([]string, error) {

	w.widget.view.Frame = false

	x, _ := gui.Size()

	buf := []string{}

	buf = append(buf,
		fmt.Sprintf("\033[37;7m \033[3%sm Setup \033[37;7m \033[3%sm Result \033[37;7m \033[3%sm Edit \033[37;7m%s", w.color(tabSetup), w.color(tabResult), w.color(tabEdit), space("", 0, x-24)),
	)

	return buf, nil
}

func (w *keybarWidget) keybinding() error {

	return nil
}

func (w *keybarWidget) color(t keybarTab) string {

	if w.tab == t {
		return "2;7"
	}

	return "7;7"
}

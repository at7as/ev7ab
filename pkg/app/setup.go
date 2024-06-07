package app

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type SetupViewWidget struct{}

func NewSetupViewWidget() *SetupViewWidget {
	return &SetupViewWidget{}
}

func (w *SetupViewWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("setup", -1, 0, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()

	v.Frame = true
	v.Title = " Setup "
	fmt.Fprintf(v, " \n")
	for i, it := range app.setup.l {
		if i == app.setup.ai {
			fmt.Fprintf(v, "  \033[37;7m%s  %s\033[0m\n\n", it.title, it.value)
		} else {
			fmt.Fprintf(v, "  %s  %s\n\n", it.title, it.value)
		}
	}

	return nil
}

type itemList struct {
	l []*item
	// m map[string]*item
	a  *item
	ai int
}

func newItemList(l [][2]string) *itemList {

	il := itemList{}

	il.l = make([]*item, len(l))
	// il.m = make(map[string]*item)
	maxLen := 0
	for i, v := range l {
		il.l[i] = &item{key: v[0], title: v[0], value: v[1]}
		// il.m[v[0]] = il.l[i]
		maxLen = max(maxLen, len(v[0]))
	}
	for i, v := range il.l {
		il.l[i].title = v.title + strings.Repeat(" ", maxLen-len(v.title))
	}

	il.assignActive()

	return &il
}

func (il *itemList) assignActive() {

	if len(il.l) > 0 {
		il.a = il.l[il.ai]
	}

}

func (il *itemList) selectUp() {

	if il.ai > 0 {
		il.ai--
	}
	il.assignActive()

}

func (il *itemList) selectDown() {

	if il.ai < len(il.l)-1 {
		il.ai++
	}
	il.assignActive()

}

type item struct {
	key   string
	title string
	value string
}

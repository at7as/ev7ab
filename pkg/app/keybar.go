package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type keybarWidget struct{}

func newKeybarWidget() *keybarWidget {
	return &keybarWidget{}
}

func (w *keybarWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("keybar", -1, maxY-2, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	v.Frame = false

	fmt.Fprintf(v, " \033[37;7mF1\033[0m Help  \033[37;7mF2\033[0m\033[3%s;%sm Setup \033[0m \033[37;7mF3\033[0m\033[3%s;%sm Result \033[0m \033[37;7mF4\033[0m\033[3%s;%sm Edit \033[0m", w.getColorF(0), w.getColorB(0), w.getColorF(1), w.getColorB(1), w.getColorF(2), w.getColorB(2))
	// switch app.view {
	// case 0:
	// 	fmt.Fprintf(v, " %s\033[37;7m^S\033[0m Save  \033[37;7m^L\033[0m Load ", strings.Repeat(" ", max(0, maxX-41-16)))
	// case 1:
	// 	fmt.Fprintf(v, " %s\033[37;7m↑↓\033[0m Select  \033[37;7mSpace\033[0m Stop  \033[37;7mEnter\033[0m Run ", strings.Repeat(" ", max(0, maxX-41-31)))
	// case 2:
	// 	fmt.Fprintf(v, " %s\033[37;7m↑↓\033[0m Select \033[37;7m↓ \033[0m Down ", strings.Repeat(" ", max(0, maxX-41-14)))
	// }

	return nil
}

func (w *keybarWidget) getColorF(i int) string {

	if app.view == i {
		return "2"
	}

	return "7"
}

func (w *keybarWidget) getColorB(i int) string {

	if app.view == i {
		return "7"
	}

	return "4"
}

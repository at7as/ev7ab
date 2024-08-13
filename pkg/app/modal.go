package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type helpBox struct {
	*widget
}

func newHelpBox() *widget {

	box := &helpBox{}
	box.widget = newWidget(box, "help")

	return box.widget
}

func (c *helpBox) transform(x int, y int) (int, int, int, int) {

	return 6, 3, x - 7, y - 4
}

func (c *helpBox) render() ([]string, error) {

	c.widget.view.Frame = true
	c.widget.view.Title = " Help "

	t := ""
	for i := 0; i <= 7; i++ {
		for _, j := range []int{1, 4, 7} {
			t += fmt.Sprintf("\033[3%d;%dm#\033[0m", i, j)
		}
	}

	return []string{t}, nil
}

func (c *helpBox) keybinding() error {

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEsc, gocui.ModNone, c.close); err != nil {
		return err
	}

	return nil
}

func (c *helpBox) close(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.setTabCurrent(app.s.keybar.getTab()); err != nil {
		return err
	}

	return app.closeModal()
}

type setupItemBox struct {
	*widget
	key   string
	value string
}

func newSetupItemBox(key, value string) *widget {

	box := &setupItemBox{key: key, value: value}
	box.widget = newWidget(box, "setupitem")

	return box.widget
}

func (c *setupItemBox) transform(x int, y int) (int, int, int, int) {

	return x/2 - 30, y/2 - 2, x/2 + 30, y / 2
}

func (c *setupItemBox) render() ([]string, error) {

	gui.Cursor = true

	c.widget.view.Editable = true
	c.widget.view.Frame = true
	c.widget.view.Title = fmt.Sprintf(" %v ", c.key)

	buf := []string{}

	buf = append(buf,
		fmt.Sprintf("%v", c.value),
	)

	x, _ := c.widget.view.Size()

	if err = c.widget.view.SetCursor(min(len(c.value), x-1), 0); err != nil {
		return buf, err
	}

	return buf, nil
}

func (c *setupItemBox) keybinding() error {

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEsc, gocui.ModNone, c.close); err != nil {
		return err
	}

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEnter, gocui.ModNone, c.enter); err != nil {
		return err
	}

	return nil
}

func (c *setupItemBox) close(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.setCurrent("setup"); err != nil {
		return err
	}

	return app.closeModal()
}

func (c *setupItemBox) enter(_ *gocui.Gui, _ *gocui.View) error {

	value, _ := c.widget.view.Line(0)

	app.s.setup.setList(c.key, value)

	// update app lab config

	return c.close(nil, nil)
}

type nodeSizeBox struct{}

func (c *nodeSizeBox) transform(x int, y int) (int, int, int, int) {

	return 0, 0, x, y
}

func (c *nodeSizeBox) render() ([]string, error) {

	return []string{}, nil
}

func (c *nodeSizeBox) keybinding() error {

	return nil
}

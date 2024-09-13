package app

import (
	"fmt"
	"strconv"
	"strings"

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

	buf := []string{}
	buf = append(buf, " ")
	buf = append(buf, fmt.Sprintf(" %s Show/hide Help", c.button("F1")))
	buf = append(buf, fmt.Sprintf(" %s Show Setup", c.button("F2")))
	buf = append(buf, fmt.Sprintf(" %s Show Result", c.button("F3")))
	buf = append(buf, fmt.Sprintf(" %s Show Edit", c.button("F4")))
	buf = append(buf, fmt.Sprintf(" %s Quit", c.button("^Q")))
	buf = append(buf, " ")

	switch app.v.keybar.tab {

	case tabSetup:
		buf = append(buf, fmt.Sprintf(" %s Change value", c.button("Enter")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s Move down/up", c.button("↓↑")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s Save lab data", c.button("^S")))
		buf = append(buf, fmt.Sprintf(" %s Load lab data", c.button("^L")))

	case tabResult:
		buf = append(buf, fmt.Sprintf(" %s Toggle run state", c.button("Enter")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s Move down/up     %s Show/hide inactive projects", c.button("↓↑"), c.button("Tab")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s New project      %s Hold project", c.button("^N"), c.button("^H")))
		buf = append(buf, fmt.Sprintf(" %s Edit project     %s Activate project", c.button("^E"), c.button("^A")))
		buf = append(buf, fmt.Sprintf(" %s Dub project      %s Terminate project", c.button("^D"), c.button("^T")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s  Select/deselect project          %s Save lab data", c.button("Space"), c.button("^S")))
		buf = append(buf, fmt.Sprintf(" %s Select/deselect all projects     %s Load lab data", c.button("^Space"), c.button("^L")))

	case tabEdit:
		buf = append(buf, fmt.Sprintf(" %s New project          %s Insert stage", c.button("^N"), c.button("Insert")))
		buf = append(buf, fmt.Sprintf(" %s Edit project         %s Delete stage", c.button("^E"), c.button("Delete")))
		buf = append(buf, fmt.Sprintf(" %s Dub project          %s  Insert node", c.button("^D"), c.button("Enter")))
		buf = append(buf, fmt.Sprintf(" %s Save project         %s Edit node", c.button("^S"), c.button("^Space")))
		buf = append(buf, fmt.Sprintf(" %s Validate project     %s     Delete node", c.button("^V"), c.button("<x")))
		buf = append(buf, fmt.Sprintf("                         %s  Edit link", c.button("Space")))
		buf = append(buf, " ")
		buf = append(buf, fmt.Sprintf(" %s  Cancel links/edits", c.button("Esc")))
		buf = append(buf, fmt.Sprintf(" %s Move cursor/link down/up/left/right", c.button("↓↑←→")))

	}

	return buf, nil
}

func (c *helpBox) keybinding() error {

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEsc, gocui.ModNone, c.close); err != nil {
		return err
	}

	return nil
}

func (c *helpBox) close(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.setTabCurrent(app.v.keybar.tab); err != nil {
		return err
	}

	return app.closeModal()
}

func (c *helpBox) button(caption string) string {

	return fmt.Sprintf("\033[37;7m%s\033[0m", caption)
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
		fmt.Sprint(c.value),
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

	app.v.setup.setListValue(c.key, value)

	switch c.key {
	case "Size":
		cfg := app.s.lab.GetConfig()
		size, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		cfg.Size = size
		if err = app.s.lab.Setup(cfg); err != nil {
			return err
		}
	case "Aggr":
		cfg := app.s.lab.GetConfig()
		cfg.Aggr = value
		if err = app.s.lab.Setup(cfg); err != nil {
			return err
		}
	case "Proc":
		cfg := app.s.lab.GetConfig()
		cfg.Proc = value
		if err = app.s.lab.Setup(cfg); err != nil {
			return err
		}
	case "Goal":
		cfg := app.s.lab.GetConfig()
		goal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		cfg.Goal = goal
		if err = app.s.lab.Setup(cfg); err != nil {
			return err
		}
	case "Duel":
		cfg := app.s.lab.GetConfig()
		duel, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		cfg.Duel = duel
		if err = app.s.lab.Setup(cfg); err != nil {
			return err
		}
	case "LabFile":
		app.s.setup.set(c.key, value)
		app.c.labfile = value
	case "In":
		app.s.setup.set(c.key, value)
		app.c.in = []int{}
		for _, str := range strings.Split(value, " ") {
			v, _ := strconv.Atoi(str)
			app.c.in = append(app.c.in, v)
		}
	case "Out":
		app.s.setup.set(c.key, value)
		app.c.out, _ = strconv.Atoi(value)
	default:
		if err = app.s.prod.Setup(c.key, value); err != nil {
			return err
		}
	}

	return c.close(nil, nil)
}

type nodeSizeBox struct {
	*widget
}

func newNodeSizeBox() *widget {

	box := &nodeSizeBox{}
	box.widget = newWidget(box, "nodesize")

	return box.widget
}

func (c *nodeSizeBox) transform(x int, y int) (int, int, int, int) {

	return x/2 - 30, y/2 - 2, x/2 + 30, y / 2
}

func (c *nodeSizeBox) render() ([]string, error) {

	gui.Cursor = true

	c.widget.view.Editable = true
	c.widget.view.Frame = true
	c.widget.view.Title = " Size "

	buf := []string{}

	value := fmt.Sprint(app.v.edit.getNodeSize())
	buf = append(buf,
		fmt.Sprint(value),
	)

	x, _ := c.widget.view.Size()

	if err = c.widget.view.SetCursor(min(len(value), x-1), 0); err != nil {
		return buf, err
	}

	return buf, nil
}

func (c *nodeSizeBox) keybinding() error {

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEsc, gocui.ModNone, c.close); err != nil {
		return err
	}

	if err = gui.SetKeybinding(c.widget.name, gocui.KeyEnter, gocui.ModNone, c.enter); err != nil {
		return err
	}

	return nil
}

func (c *nodeSizeBox) close(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.setCurrent("edit"); err != nil {
		return err
	}

	return app.closeModal()
}

func (c *nodeSizeBox) enter(_ *gocui.Gui, _ *gocui.View) error {

	str, _ := c.widget.view.Line(0)
	value, _ := strconv.Atoi(str)

	app.v.edit.setNodeSize(value)
	app.v.edit.draft.measure()
	app.v.edit.mark()

	return c.close(nil, nil)
}

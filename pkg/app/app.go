package app

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui
var app *application
var err error

type application struct {
	c config
	s state
	v view
}

type config struct {
	cfgfile, labfile string
	in, out          int
}

type state struct {
	status applicationStatus
	lab    *lab.Lab
	prod   lab.Producer
	setup  *dict
	ev     []*project
}

type view struct {
	modal     *widget
	keybar    *keybarWidget
	statusbar *statusbarWidget
	setup     *setupWidget
	result    *resultWidget
	edit      *editWidget
}

type applicationStatus int

const (
	appIdle applicationStatus = iota
	appRun
	appWait
)

func createApplication(prod lab.Producer, cfgFile string, debug bool) {

	app = &application{}

	app.c = config{
		cfgfile: cfgFile,
		labfile: "./ev.lab",
		in:      2,
		out:     2,
	}

	app.s = state{
		status: appIdle,
		prod:   prod,
		lab:    lab.New(prod, debug),
		setup:  newDict(defaultSetup),
		ev:     make([]*project, 0),
	}

	app.v = view{
		keybar:    newKeybarWidget(),
		statusbar: newStatusbarWidget(),
		setup:     newSetupWidget(),
		result:    newResultWidget(),
		edit:      newEditWidget(),
	}

}

func (a *application) keybinding() error {

	if err = gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, a.quit); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, a.toggleHelp); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF2, gocui.ModNone, a.showSetup); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF3, gocui.ModNone, a.showResult); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF4, gocui.ModNone, a.showEdit); err != nil {
		return err
	}

	if err = a.v.keybar.keybinding(); err != nil {
		return err
	}

	if err = a.v.statusbar.keybinding(); err != nil {
		return err
	}

	if err = a.v.setup.keybinding(); err != nil {
		return err
	}

	if err = a.v.result.keybinding(); err != nil {
		return err
	}

	if err = a.v.edit.keybinding(); err != nil {
		return err
	}

	return nil
}

func (a *application) quit(_ *gocui.Gui, _ *gocui.View) error {

	if err = a.s.lab.Close(); err != nil {
		return err
	}

	return gocui.ErrQuit
}

func (a *application) toggleHelp(_ *gocui.Gui, _ *gocui.View) error {

	if a.v.modal == nil {

		if err = a.openModal(newHelpBox()); err != nil {
			return err
		}

	} else {

		open := a.v.modal.name != "help"

		if err = a.setTabCurrent(a.v.keybar.tab); err != nil {
			return err
		}

		if err = a.closeModal(); err != nil {
			return err
		}

		if open {
			if err = a.openModal(newHelpBox()); err != nil {
				return err
			}
		}

	}

	return nil
}

func (a *application) showSetup(_ *gocui.Gui, _ *gocui.View) error {

	if err = a.closeModal(); err != nil {
		return err
	}
	a.v.keybar.setTab(tabSetup)

	return a.setTabCurrent(tabSetup)
}

func (a *application) showResult(_ *gocui.Gui, _ *gocui.View) error {

	if err = a.closeModal(); err != nil {
		return err
	}
	a.v.keybar.setTab(tabResult)

	return a.setTabCurrent(tabResult)
}

func (a *application) showEdit(_ *gocui.Gui, _ *gocui.View) error {

	if err = a.closeModal(); err != nil {
		return err
	}

	if app.v.keybar.tab == tabResult {
		p := app.v.result.getFocused()
		if p != nil {
			if app.v.edit.draft == nil {
				app.v.edit.setDraft(p)
			} else {
				if app.v.edit.draft.status != projectEdit {
					app.v.edit.setDraft(p)
				}
			}
		}
	}

	a.v.keybar.setTab(tabEdit)

	return a.setTabCurrent(tabEdit)
}

func (a *application) idle() bool {

	return a.s.status == appIdle
}

func (a *application) setCurrent(name string) error {

	if _, err = gui.SetCurrentView(name); err != nil {
		return err
	}
	if _, err = gui.SetViewOnTop(name); err != nil {
		return err
	}

	return nil
}

func (a *application) setTabCurrent(tab keybarTab) error {

	if _, err = gui.SetCurrentView(string(tab)); err != nil {
		return err
	}
	if _, err = gui.SetViewOnTop(string(tab)); err != nil {
		return err
	}

	return nil
}

func (a *application) openModal(w *widget) error {

	if a.v.modal != nil {
		a.v.modal.clean()
	}
	a.v.modal = w

	if err = a.v.modal.Layout(gui); err != nil {
		return err
	}

	if err = a.v.modal.keybinding(); err != nil {
		return err
	}

	if err = a.setCurrent(w.name); err != nil {
		return err
	}

	return nil
}

func (a *application) closeModal() error {

	gui.Cursor = false

	if a.v.modal != nil {
		if err = a.v.modal.clean(); err != nil {
			return err
		}
		a.v.modal = nil
	}

	return nil
}

func (a *application) loadConfig() error {

	f, err := os.Open(a.c.cfgfile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var cfg map[string]json.RawMessage
	if err = json.Unmarshal(b, &cfg); err != nil {
		return err
	}

	for k, v := range cfg {
		if k == "app" {
			var obj map[string]string
			if err = json.Unmarshal(v, &obj); err != nil {
				return err
			}
			for kk, vv := range obj {
				switch kk {
				case "LabFile":
					a.s.setup.set(kk, vv)
					a.c.labfile = vv
				case "In":
					a.s.setup.set(kk, vv)
					a.c.in, _ = strconv.Atoi(vv)
				case "Out":
					a.s.setup.set(kk, vv)
					a.c.out, _ = strconv.Atoi(vv)
				}
			}
		}
		if k == "prod" {
			var obj map[string]string
			if err = json.Unmarshal(v, &obj); err != nil {
				return err
			}
			for kk, vv := range obj {
				a.s.setup.set(kk, vv)
			}
		}
	}

	return nil
}

func (a *application) loadProd() error {

	d := make([]string, len(defaultSetup))
	for i, v := range defaultSetup {
		d[i] = v.key
	}

	s := make(map[string]string)
	for _, item := range a.s.setup.l {
		if !slices.Contains(d, item.key) {
			s[item.key] = item.value
		}
	}

	return a.s.prod.Load(s)
}

func (a *application) saveLab() {

	b, err := a.s.lab.Export()
	if err != nil {
		a.wait(err.Error())
		return
	}

	buf := bytes.Buffer{}
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		a.wait(err.Error())
		return
	}
	w.Write(b)
	w.Close()

	if err = os.WriteFile(a.c.labfile, buf.Bytes(), 0644); err != nil {
		a.wait(err.Error())
		return
	}

	a.apply(appIdle, "")
	a.v.statusbar.mark()
	a.update(nil)

}

func (a *application) loadLab() {

	f, err := os.Open(a.c.labfile)
	if err != nil {
		a.wait(err.Error())
		return
	}
	defer f.Close()

	d, err := zlib.NewReader(f)
	if err != nil {
		a.wait(err.Error())
		return
	}
	defer d.Close()

	b, err := io.ReadAll(d)
	if err != nil {
		a.wait(err.Error())
		return
	}

	if err = a.s.lab.Import(b); err != nil {
		a.wait(err.Error())
		return
	}

	labcfg := a.s.lab.GetConfig()
	a.s.setup.set("Size", strconv.Itoa(labcfg.Size))
	a.s.setup.set("Aggr", labcfg.Aggr)
	a.s.setup.set("Proc", labcfg.Proc)
	a.s.setup.set("Goal", strconv.FormatBool(labcfg.Goal))
	a.s.setup.set("Duel", strconv.FormatBool(labcfg.Duel))

	a.importProjects()

	a.apply(appIdle, "")
	a.v.statusbar.mark()
	a.v.setup.mark()
	a.v.result.mark()
	a.v.edit.setDraft(nil)
	a.update(nil)

}

func (a *application) importProjects() {

	list := a.s.lab.GetProjects()
	a.s.ev = make([]*project, len(list))

	for i, id := range list {

		status := projectActive
		if !a.s.lab.ProjectStatus(id) {
			status = projectHolded
		}

		p := &project{
			id:       id,
			status:   status,
			stat:     stat{},
			model:    createProjectModel(a.s.lab.ProjectLayout(id)),
			draft:    nil,
			selected: false,
		}

		p.refine()
		p.model.measure()
		p.stat.size = p.model.size
		p.stat.volume = p.model.volume

		a.s.ev[i] = p

	}

}

func (a *application) apply(s applicationStatus, t string) {

	a.s.status = s
	a.v.statusbar.setColor(s)
	a.v.statusbar.setText(t)

}

func (a *application) wait(t string) {

	a.s.status = appWait
	a.v.statusbar.setColor(appWait)
	a.v.statusbar.setText(t)
	a.update(nil)

	go func() {
		time.Sleep(2 * time.Second)
		a.apply(appIdle, "")
		a.update(nil)
	}()

}

func (a *application) update(f func(g *gocui.Gui) error) {

	if f != nil {
		gui.Update(f)
	} else {
		gui.Update(func(g *gocui.Gui) error {
			return nil
		})
	}

}

// Run ...
func Run(prod lab.Producer, appConfigFile string, debug bool) {

	if gui, err = gocui.NewGui(gocui.OutputNormal); err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	createApplication(prod, appConfigFile, debug)
	if err = app.loadConfig(); err != nil {
		log.Panicln(err)
	}
	if err = app.loadProd(); err != nil {
		log.Panicln(err)
	}

	gui.SetManager(
		app.v.keybar.widget,
		app.v.statusbar.widget,
		app.v.setup.widget,
		app.v.result.widget,
		app.v.edit.widget,
	)

	if err = app.keybinding(); err != nil {
		log.Panicln(err)
	}

	gui.InputEsc = true

	app.update(func(g *gocui.Gui) error {
		return app.showSetup(nil, nil)
	})

	if err = gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

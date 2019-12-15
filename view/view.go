package view

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/wolffcm/fluxcui"
)

type Config struct {
	EditorText string
}

type cui struct {
	cfg *Config
	m fluxcui.Model
	c fluxcui.Controller

	lg *linegraph
}

func NewView(cfg *Config, m fluxcui.Model, c fluxcui.Controller) fluxcui.View {
	return &cui{
		cfg: cfg,
		m: m,
		c: c,
		lg: newLinegraph(),
	}
}

func (c *cui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.Mouse = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(c.layout)

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}
	if err := c.setControlPanelKeybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (c *cui) nextView(g *gocui.Gui, v *gocui.View) error {
	var newView string
	switch n := v.Name(); n {
	case "control": newView = "editor"
	case "editor": newView = "control"
	}
	if _, err := setCurrentViewOnTop(g, newView); err != nil {
		return err
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func (c *cui) layout(g *gocui.Gui) error {
	if _, err := g.View("errors"); err == nil {
		mustWriteMessage(g, "calling layout...")
	}
	maxX, maxY := g.Size()

	row1y := maxY - int(float64(maxY) * .2)
	if maxY -row1y < 12 {
		row1y = maxY - 12
	}

	linegraph, err := g.SetView("linegraph", 0, 0, maxX - 1, row1y- 1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		linegraph.Title = "FluxCUI"
	}

	edX := 10

	control, err := g.SetView("control", 0, row1y, edX - 1, maxY - 1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		control.Clear()
		if _, err := fmt.Fprint(control, "Run"); err != nil {
			return err
		}
	}

	errPanelX := maxX - (maxX / 3)

	editor, err := g.SetView("editor", edX, row1y, errPanelX - 1, maxY - 1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editor.Title = "editor"
		editor.Editable = true
		editor.Wrap = true
		if _, err = fmt.Fprintf(editor, c.cfg.EditorText); err != nil {
			return err
		}

		if _, err := setCurrentViewOnTop(g, "editor"); err != nil {
			return err
		}
	}

	if errPanel, err := g.SetView("errors", errPanelX, row1y, maxX - 1, maxY - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		errPanel.Title = "errors"
		errPanel.Autoscroll = true
	}

	if err := c.lg.update(g, c.m); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}


func writeError(g *gocui.Gui, userErr error) error {
	return writeMessage(g, userErr.Error())
}

func writeMessage(g *gocui.Gui, msg string) error {
	v, err := g.View("errors")
	if err != nil {
		return err
	}
	ts := time.Now().Format(time.Stamp)
	if _, err := fmt.Fprintf(v, "%v: %v\n", ts, msg); err != nil {
		return err
	}
	return nil
}

func mustWriteMessage(g *gocui.Gui, msg string) {
	if err := writeMessage(g, msg); err != nil {
		panic(err)
	}
}

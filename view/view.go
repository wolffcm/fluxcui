package view

import (
	"github.com/jroimartin/gocui"
	"github.com/wolffcm/fluxcui"
)

type Config struct {
	EditorText string
}

type cui struct {
	cfg *Config
	m   fluxcui.Model
	c   fluxcui.Controller

	lg *lineGraph
}

func NewView(cfg *Config, m fluxcui.Model, c fluxcui.Controller) fluxcui.View {
	return &cui{
		cfg: cfg,
		m:   m,
		c:   c,
		lg:  newLineGraph(),
	}
}

func (c *cui) Run() error {
	g, err := gocui.NewGui(gocui.Output256)
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
	if err := c.setEditorKeybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

const (
	controlView = "control"
	editorView  = "editor"
	logView     = "log"
	graphView   = "graph"
)

func (c *cui) layout(g *gocui.Gui) error {
	defer func() {
		if p := recover(); p != nil {
			mustWriteMessagef(g, "panic: %v", p)
		}
	}()

	maxX, maxY := g.Size()
	row1y := maxY - int(float64(maxY)*.2)
	if maxY-row1y < 12 {
		row1y = maxY - 12
	}

	if err := doGraphView(g, 0, 0, maxX-1, row1y-1); err != nil {
		return err
	}

	edX := 10
	if err := doControlView(g, 0, row1y, edX-1, maxY-1); err != nil {
		return err
	}

	errPanelX := maxX - (maxX / 3)

	if err := c.doEditorView(g, edX, row1y, errPanelX-1, maxY-1); err != nil {
		return err
	}

	if err := c.doLogView(g, errPanelX, row1y, maxX-1, maxY-1); err != nil {
		return err
	}

	if err := c.lg.update(g, c.m); err != nil {
		return err
	}
	return nil
}

func (c *cui) nextView(g *gocui.Gui, v *gocui.View) error {
	var newView string
	switch n := v.Name(); n {
	case controlView:
		newView = editorView
	case editorView:
		newView = controlView
	}
	if _, err := setCurrentViewOnTop(g, newView); err != nil {
		return err
	}
	return nil
}

func setcurrentView(g *gocui.Gui, v *gocui.View) error {
	_, err := setCurrentViewOnTop(g, v.Name())
	return err
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

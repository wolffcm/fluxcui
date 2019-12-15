package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (c *cui) setControlPanelKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("control", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("control", gocui.MouseLeft, gocui.ModNone, c.controlMouseClick); err != nil {
		return err
	}
	if err := g.SetKeybinding("control", gocui.KeyEnter, gocui.ModNone, c.controlMouseClick); err != nil {
		return err
	}

	return nil
}

func doControlView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	control, err := g.SetView("control", x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		control.Clear()
		if _, err := fmt.Fprint(control, "Run"); err != nil {
			return err
		}
	}
	return nil
}

func (c *cui) controlMouseClick(g *gocui.Gui, v *gocui.View) error {
	var err error
	if _, err = setCurrentViewOnTop(g, v.Name()); err != nil {
		return err
	}
	_, cy := v.Cursor()
	var line string
	if line, err = v.Line(cy); err != nil {
		line = ""
	}

	switch line {
	case "Run":
		if err := writeMessage(g, "executing query..."); err != nil {
			return err
		}
		ev, err := g.View("editor")
		if err != nil {
			return err
		}
		q := ev.Buffer()
		if err := c.c.Query(q); err != nil {
			if err := writeError(g, err); err != nil {
				return err
			}
		}
	}

	return nil
}

package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (c *cui) setControlPanelKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(controlView, 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(controlView, gocui.MouseLeft, gocui.ModNone, c.doMenuItem); err != nil {
		return err
	}
	if err := g.SetKeybinding(controlView, gocui.KeyEnter, gocui.ModNone, c.doMenuItem); err != nil {
		return err
	}
	if err := g.SetKeybinding(controlView, gocui.KeyArrowUp, gocui.ModNone, c.doCursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(controlView, gocui.KeyArrowDown, gocui.ModNone, c.doCursorDown); err != nil {
		return err
	}

	return nil
}

type menuAction func(g *gocui.Gui, c *cui) error

type menuItem struct {
	name   string
	action menuAction
}

var menuItems = []menuItem{
	{
		name:   "Run",
		action: runQuery,
	},
	{
		name:   "Clear",
		action: clear,
	},
}

func doControlView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	control, err := g.SetView(controlView, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		control.Highlight = true
		control.Title = "menu"
		needNewline := false
		for _, item := range menuItems {
			if needNewline {
				if _, err := fmt.Fprint(control, "\n"); err != nil {
					return err
				}
			} else {
				needNewline = true
			}
			if _, err := fmt.Fprint(control, item.name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *cui) doMenuItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	if _, err = setCurrentViewOnTop(g, v.Name()); err != nil {
		return err
	}
	_, cy := v.Cursor()
	var line string
	if line, err = v.Line(cy); err != nil {
		line = ""
	}

	if f := getAction(menuItems, line); f != nil {
		return f(g, c)
	}
	return nil
}

func (c *cui) doCursorUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if y > 0 {
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}

	return nil
}

func (c *cui) doCursorDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if y < len(menuItems)-1 {
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}

	return nil
}

func getAction(items []menuItem, name string) menuAction {
	for _, i := range items {
		if i.name == name {
			return i.action
		}
	}
	return nil
}

func runQuery(g *gocui.Gui, c *cui) error {
	if err := writeMessage(g, "executing query"); err != nil {
		return err
	}
	ev, err := g.View(editorView)
	if err != nil {
		return err
	}
	q := ev.Buffer()
	if err := c.c.Query(q); err != nil {
		if err := writeError(g, err); err != nil {
			return err
		}
	}
	return nil
}

func clear(g *gocui.Gui, c *cui) error {
	if err := writeMessage(g, "clearing"); err != nil {
		return err
	}

	v, err := g.View(editorView)
	if err != nil {
		return err
	}
	v.Clear()

	v, err = g.View(graphView)
	if err != nil {
		return err
	}
	v.Clear()

	return nil
}

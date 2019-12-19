package view

import (
	"github.com/jroimartin/gocui"
)

func (c *cui) doLogView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView(logView, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "log"
		v.Autoscroll = true
	}
	return nil
}

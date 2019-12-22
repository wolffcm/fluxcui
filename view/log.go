package view

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

func (c *cui) doLogView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView(logView, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "log"
		v.Autoscroll = true
		v.Wrap = true
	}
	return nil
}

func writeError(g *gocui.Gui, userErr error) error {
	return writeMessage(g, userErr.Error())
}

func writeMessage(g *gocui.Gui, msg string) error {
	v, err := g.View(logView)
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

func mustWriteMessagef(g *gocui.Gui, format string, args ...interface{}) {
	msg := fmt.Sprintf(fmt.Sprintf(format, args...))
	if err := writeMessage(g, msg); err != nil {
		panic(err)
	}
}

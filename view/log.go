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

func logError(g *gocui.Gui, userErr error) error {
	return logMessage(g, userErr.Error())
}

func logMessage(g *gocui.Gui, msg string) error {
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

func mustLogMessage(g *gocui.Gui, msg string) {
	if err := logMessage(g, msg); err != nil {
		panic(err)
	}
}

func mustLogMessagef(g *gocui.Gui, format string, args ...interface{}) {
	msg := fmt.Sprintf(fmt.Sprintf(format, args...))
	if err := logMessage(g, msg); err != nil {
		panic(err)
	}
}

func (c *cui) logVerbose(g *gocui.Gui, msg string) error {
	if c.cfg.Verbose {
		return logMessage(g, msg)
	}
	return nil
}

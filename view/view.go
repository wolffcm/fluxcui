package view

import (
	"fmt"
	"math"
	"time"

	"github.com/exrook/drawille-go"
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
	cv drawille.Canvas
}

func NewView(cfg *Config, m fluxcui.Model, c fluxcui.Controller) fluxcui.View {
	return &cui{
		cfg: cfg,
		m: m,
		c: c,
		cv: drawille.NewCanvas(),
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

	if err := g.SetKeybinding("control", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("control", gocui.MouseLeft, gocui.ModNone, c.controlMouseClick); err != nil {
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
	writeMessage(g, "calling layout...")
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

	// TODO(cwolff): only recompute this graph if size changes
	if err := c.writeView(linegraph); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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


type canvasPoint struct {
	X, Y float64
}

func (c *cui) writeView(v *gocui.View) error {
	c.cv.Clear()

	ss := c.m.Series()
	vxd, vyd := v.Size()
	cxd, cyd := vxd*2, vyd*4
	xformer := getTransformer(ss, float64(cxd), float64(cyd))
	for _, s := range ss {
		tps := s.Data
		oldPt := xformer(tps[0])
		for _, tp := range tps[1:] {
			cp := xformer(tp)
			c.cv.DrawLine(oldPt.X, oldPt.Y, cp.X, cp.Y)
			oldPt = cp
		}
	}

	v.Clear()
	if _, err := v.Write([]byte(c.cv.String())); err != nil {
		return err
	}

	return nil
}

type pointTransformer func(point fluxcui.TimePoint) canvasPoint

func getTransformer(ss []fluxcui.Series, cxd, cyd float64) pointTransformer {
	cxd--
	cyd--

	minT, maxT := int64(math.MaxInt64), int64(math.MinInt64)
	for _, s := range ss {
		nPts := len(s.Data)
		if s.Data[0].T.UnixNano() < minT {
			minT = s.Data[0].T.UnixNano()
		}
		if s.Data[nPts - 1].T.UnixNano() > maxT {
			maxT = s.Data[nPts - 1].T.UnixNano()
		}
	}

	txLen := maxT - minT
	xScale := cxd / float64(txLen)
	xTranslate := -minT

	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	for _, s := range ss {
		tps := s.Data
		for _, tp := range tps {
			minY = math.Min(minY, tp.V)
			maxY = math.Max(maxY, tp.V)
		}
	}

	tyLen := maxY - minY
	yScale := cyd / tyLen
	yTranslate := -minY

	return func(tp fluxcui.TimePoint) canvasPoint {
		t := tp.T.UnixNano()
		return canvasPoint{
			X: float64(t+xTranslate) * xScale,
			Y: cyd - ((tp.V + yTranslate) * yScale),
		}
	}
}

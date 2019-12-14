package view

import (
	"math"

	"github.com/exrook/drawille-go"
	"github.com/jroimartin/gocui"
	"github.com/wolffcm/fluxcui"
)

type cui struct {
	m fluxcui.Model
	cv drawille.Canvas
}

func NewView(m fluxcui.Model) fluxcui.View {
	return &cui{
		m: m,
		cv: drawille.NewCanvas(),
	}
}

func (c *cui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(c.layout)

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (c *cui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 2, 1, maxX - 3, maxY - 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := c.writeView(v); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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

func scaleToCanvas(tps []fluxcui.TimePoint, cxd, cyd float64) []canvasPoint {
	cxd--
	cyd--

	// Assume data sorted by time.
	numPts := len(tps)
	minT := tps[0].T.UnixNano()
	maxT := tps[numPts - 1].T.UnixNano()

	txLen := maxT - minT
	xScale := cxd / float64(txLen)
	xTranslate := -minT

	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	for _, tp := range tps {
		minY = math.Min(minY, tp.V)
		maxY = math.Max(maxY, tp.V)
	}

	tyLen := maxY - minY
	yScale := cyd / tyLen
	yTranslate := -minY

	cps := make([]canvasPoint, numPts)
	for i, tp := range tps {
		t := tp.T.UnixNano()
		cps[i] = canvasPoint{
			X: float64(t+xTranslate) * xScale,
			Y: cyd - ((tp.V + yTranslate) * yScale),
		}
	}
	return cps
}

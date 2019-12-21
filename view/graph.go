package view

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/wolffcm/drawille-go"
	"github.com/wolffcm/fluxcui"
	"io"
	"math"
	"time"
)

func doGraphView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(graphView, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Title = "FluxCUI"
	}
	return nil
}

type canvasPoint struct {
	X, Y float64
}

type lineGraph struct {
	vxSize, vySize int
	ts             time.Time

	canvas drawille.Canvas
}

func newLineGraph() *lineGraph {
	return &lineGraph{
		canvas: drawille.NewCanvas(drawille.SaturateOnOverwrite()),
	}
}

func (lg *lineGraph) update(g *gocui.Gui, m fluxcui.Model) error {
	v, err := g.View(graphView)
	if err != nil {
		return err
	}

	vxd, vyd := v.Size()
	if lg.vxSize == vxd && lg.vySize == vyd && !m.Timestamp().After(lg.ts) {
		mustWriteMessage(g, "skipping update")
		return nil
	} else {
		mustWriteMessage(g, "updating")
	}

	cxd, cyd := vxd*2, vyd*4
	v.Clear()
	if err := lg.render(m, cxd, cyd, v); err != nil {
		return err
	}
	lg.ts = m.Timestamp()
	lg.vxSize = vxd
	lg.vySize = vyd

	return nil

}

func (lg *lineGraph) render(m fluxcui.Model, xd, yd int, w io.Writer) error {
	ss := m.Series()
	palette := getPalette(len(ss))
	lg.canvas.Clear(drawille.SetPalette(palette))

	translator := getTranslator(ss, float64(xd), float64(yd))
	for i, s := range ss {
		tps := s.Data
		oldPt := translator(tps[0])
		for _, tp := range tps[1:] {
			cp := translator(tp)
			lg.canvas.DrawLine(oldPt.X, oldPt.Y, cp.X, cp.Y, i)
			oldPt = cp
		}
	}

	if _, err := fmt.Fprint(w, lg.canvas.String()); err != nil {
		return err
	}

	return nil
}

type pointTranslator func(point fluxcui.TimePoint) canvasPoint

func getTranslator(ss []fluxcui.Series, cxd, cyd float64) pointTranslator {
	cxd--
	cyd--

	minT, maxT := int64(math.MaxInt64), int64(math.MinInt64)
	for _, s := range ss {
		nPts := len(s.Data)
		if s.Data[0].T.UnixNano() < minT {
			minT = s.Data[0].T.UnixNano()
		}
		if s.Data[nPts-1].T.UnixNano() > maxT {
			maxT = s.Data[nPts-1].T.UnixNano()
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

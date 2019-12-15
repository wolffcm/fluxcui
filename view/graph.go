package view

import (
	"math"
	"time"

	"github.com/exrook/drawille-go"
	"github.com/jroimartin/gocui"
	"github.com/wolffcm/fluxcui"
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

type linegraph struct {
	vxsize, vysize int
	ts             time.Time

	canvas drawille.Canvas
}

func newLinegraph() *linegraph {
	return &linegraph{
		canvas: drawille.NewCanvas(),
	}
}

func (lg *linegraph) update(g *gocui.Gui, m fluxcui.Model) error {
	v, err := g.View(graphView)
	if err != nil {
		return err
	}

	vxd, vyd := v.Size()
	if lg.vxsize == vxd && lg.vysize == vyd && !m.Timestamp().After(lg.ts) {
		mustWriteMessage(g, "skipping update (nothing to do)")
		return nil
	} else {
		mustWriteMessage(g, "updating")
	}

	lg.canvas.Clear()

	ss := m.Series()
	cxd, cyd := vxd*2, vyd*4
	xformer := getTransformer(ss, float64(cxd), float64(cyd))
	for _, s := range ss {
		tps := s.Data
		oldPt := xformer(tps[0])
		for _, tp := range tps[1:] {
			cp := xformer(tp)
			lg.canvas.DrawLine(oldPt.X, oldPt.Y, cp.X, cp.Y)
			oldPt = cp
		}
	}

	v.Clear()
	if _, err := v.Write([]byte(lg.canvas.String())); err != nil {
		return err
	}

	lg.ts = m.Timestamp()
	lg.vxsize = vxd
	lg.vysize = vyd

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

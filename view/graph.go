package view

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/wolffcm/drawille-go"
	"github.com/wolffcm/fluxcui"
	"io"
	"math"
	"sort"
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
	cui *cui

	vxSize, vySize int
	ts             time.Time

	canvas drawille.Canvas
}

func newLineGraph(cui *cui) *lineGraph {
	return &lineGraph{
		cui:    cui,
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
		return lg.cui.logVerbose(g, "skipping update")

	} else if err := lg.cui.logVerbose(g, "updating"); err != nil {
		return err
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
	if fluxcui.IsEmpty(ss) {
		return nil
	}

	palette := getPalette(len(ss))
	lg.canvas.Clear(drawille.SetPalette(palette))

	md := getMetadata(ss)
	translator := getTranslator(ss, md, float64(xd), float64(yd))

	xTicks := chooseXTicks(md)
	for _, xt := range xTicks {
		cp0 := translator(fluxcui.TimePoint{T: xt, V: md.minV})
		cp1 := translator(fluxcui.TimePoint{T: xt, V: md.maxV})
		lg.canvas.DrawLine(cp0.X, cp0.Y, cp1.X, cp1.Y, grey)
	}

	yTicks := chooseYTicks(md)
	for _, yt := range yTicks {
		cp0 := translator(fluxcui.TimePoint{T: md.minT, V: yt})
		cp1 := translator(fluxcui.TimePoint{T: md.maxT, V: yt})
		lg.canvas.DrawLine(cp0.X, cp0.Y, cp1.X, cp1.Y, grey)
	}

	for i, s := range ss {
		tps := s.Data
		oldPt := translator(tps[0])
		for _, tp := range tps[1:] {
			cp := translator(tp)
			lg.canvas.DrawLine(oldPt.X, oldPt.Y, cp.X, cp.Y, i)
			oldPt = cp
		}
	}

	for _, x := range xTicks {
		txt := x.Format("15:04:05")
		cp := translator(fluxcui.TimePoint{T: x, V: md.minV})
		ln := len(txt)
		cp.X -= float64(ln/2) * 2
		lg.canvas.SetText(int(cp.X), int(cp.Y), txt, grey)
	}

	for _, y := range yTicks {
		txt := fmt.Sprintf("%v", y)
		cp := translator(fluxcui.TimePoint{T: md.minT, V: y})
		lg.canvas.SetText(0, int(cp.Y), txt, grey)
	}

	if _, err := fmt.Fprint(w, lg.canvas.String()); err != nil {
		return err
	}

	return nil
}

type pointTranslator func(point fluxcui.TimePoint) canvasPoint

type Metadata struct {
	minT, maxT time.Time
	minV, maxV float64
}

func getMetadata(ss []fluxcui.Series) *Metadata {
	minT, maxT := time.Unix(0, math.MaxInt64), time.Unix(0, math.MinInt64)
	minV, maxV := math.MaxFloat64, -math.MaxFloat64
	for _, s := range ss {
		tps := s.Data
		for _, tp := range tps {
			minV = math.Min(minV, tp.V)
			maxV = math.Max(maxV, tp.V)
			if tp.T.After(maxT) {
				maxT = tp.T
			}
			if tp.T.Before(minT) {
				minT = tp.T
			}
		}
	}

	return &Metadata{
		minT: minT,
		maxT: maxT,
		minV: minV,
		maxV: maxV,
	}
}

func getTranslator(ss []fluxcui.Series, md *Metadata, cxd, cyd float64) pointTranslator {
	cxd--
	cyd--

	txLen := md.maxT.Sub(md.minT)
	xScale := cxd / float64(txLen)
	xTranslate := -md.minT.UnixNano()

	tyLen := md.maxV - md.minV
	yScale := cyd / tyLen
	yTranslate := -md.minV

	return func(tp fluxcui.TimePoint) canvasPoint {
		t := tp.T.UnixNano()
		return canvasPoint{
			X: float64(t+xTranslate) * xScale,
			Y: cyd - ((tp.V + yTranslate) * yScale),
		}
	}
}

func chooseXTicks(md *Metadata) []time.Time {
	start := md.minT
	stop := md.maxT

	d := stop.Sub(start)

	ds := []time.Duration{
		time.Hour * 24 * 7,
		time.Hour * 48,
		time.Hour * 24,
		time.Hour * 6,
		time.Hour,
		time.Minute * 15,
		time.Minute * 5,
		time.Minute,
		time.Second * 15,
		time.Second * 5,
		time.Second,
		time.Millisecond,
		time.Microsecond,
		time.Nanosecond,
	}

	i := sort.Search(len(ds), func(i int) bool {
		if d > ds[i] {
			return true
		}
		return false
	})
	xTicks := make([]time.Time, 0, 5)
	if i >= len(ds) {
		// Sub ns interval??
		xTicks = append(xTicks, start)
		return xTicks
	} else {
		truncDur := ds[i]
		tick := start.Truncate(truncDur)
		for ; tick.Before(stop); tick = tick.Add(truncDur) {
			if tick.After(start) {
				xTicks = append(xTicks, tick)
			}
		}
	}

	return xTicks
}

func chooseYTicks(md *Metadata) []float64 {
	yRange := md.maxV - md.minV
	unit := math.Pow(10.0, math.Trunc(math.Log10(yRange)))
	if yRange/unit <= 2 {
		unit /= 2
	}
	yTicks := make([]float64, 0, 5)
	v := unit * math.Trunc(md.minV/unit)

	for ; v < md.maxV; v += unit {
		if v > md.minV {
			yTicks = append(yTicks, v)
		}
	}

	return yTicks
}

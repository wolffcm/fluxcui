package fluxcui

import (
	"math"
	"time"
)

type CanvasPoint struct {
	X, Y float64
}

type TimePoint struct {
	T time.Time
	V float64
}

func ScaleToCanvas(tps []TimePoint, cxd, cyd float64) []CanvasPoint {
	cxd--
	cyd--

	// Assume data sorted by time.
	numPts := len(tps)
	minT := tps[0].T.UnixNano()
	maxT := tps[numPts - 1].T.UnixNano()

	txLen := maxT - minT
	xScale := float64(cxd) / float64(txLen)
	xTranslate := -minT

	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	for _, tp := range tps {
		minY = math.Min(minY, tp.V)
		maxY = math.Max(maxY, tp.V)
	}

	tyLen := maxY - minY
	yScale := float64(cyd) / tyLen
	yTranslate := -minY

	cps := make([]CanvasPoint, numPts)
	for i, tp := range tps {
		t := tp.T.UnixNano()
		cps[i] = CanvasPoint{
			X: float64(t+xTranslate) * xScale,
			Y: cyd - ((tp.V + yTranslate) * yScale),
		}
	}
	return cps
}

func GenData() []TimePoint {
	numPts := 33
	points := make([]TimePoint, numPts)
	startTime := time.Now().Truncate(time.Minute)
	for i := 0; i < numPts; i++ {
		t := startTime.Add(time.Minute*time.Duration(i))
		v := math.Sin(float64(i) * (math.Pi/8))
		points[i] = TimePoint{T: t, V: v}
	}
	return points
}


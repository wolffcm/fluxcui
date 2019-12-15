package model

import (
	"math"
	"time"

	"github.com/wolffcm/fluxcui"
)

type sin struct {
	series []fluxcui.Series
	ts     time.Time
}

func NewModel() fluxcui.Model {
	return &sin{
		series: genData(),
		ts:     time.Now(),
	}
}

func (s *sin) Timestamp() time.Time {
	return s.ts
}

func (s *sin) Query(fluxSrc string) error {
	return nil
}

func (s *sin) Series() []fluxcui.Series {
	return s.series
}

func genData() []fluxcui.Series {
	ss := make([]fluxcui.Series, 0, 3)
	ss = append(ss, fluxcui.Series{
		Tags: map[string]string{
			"step": "pi/8",
		},
		Data: genSeries(math.Pi / 8),
	})
	ss = append(ss, fluxcui.Series{
		Tags: map[string]string{
			"step": "pi/16",
		},
		Data: genSeries(math.Pi / 16),
	})

	pts := make([]fluxcui.TimePoint, len(ss[0].Data))
	for i, pt := range ss[0].Data {
		sumPt := pt.V + ss[1].Data[i].V
		pts[i] = fluxcui.TimePoint{
			T: pt.T,
			V: sumPt,
		}
	}
	ss = append(ss, fluxcui.Series{
		Tags: map[string]string{
			"method": "sum",
		},
		Data: pts,
	})
	return ss
}

func genSeries(step float64) []fluxcui.TimePoint {
	numPts := 33
	points := make([]fluxcui.TimePoint, numPts)
	startTime := time.Now().Truncate(time.Minute)
	for i := 0; i < numPts; i++ {
		t := startTime.Add(time.Minute * time.Duration(i))
		v := math.Sin(float64(i) * step)
		points[i] = fluxcui.TimePoint{T: t, V: v}
	}
	return points
}

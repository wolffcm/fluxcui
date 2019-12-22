package mock

import (
	"fmt"
	"github.com/wolffcm/fluxcui"
	"math"
	"time"
)

type Model struct {
}

func (m Model) Query(fluxSrc string) error {
	return nil
}

func (m Model) Series() []fluxcui.Series {
	numSeries := 5
	pointsPerSeries := 80
	dur := time.Hour
	ss := make([]fluxcui.Series, numSeries)
	for i := 0; i < numSeries; i++ {
		s := fluxcui.Series{
			Tags: map[string]string{
				"v": fmt.Sprintf("%v", i),
			},
			Data: make([]fluxcui.TimePoint, pointsPerSeries),
		}
		shift := 2.0 * math.Pi * (float64(i) / float64(numSeries))
		for j := 0; j < pointsPerSeries; j++ {
			t := m.Timestamp().Add(-dur)
			increment := (dur * time.Duration(j)) / time.Duration(pointsPerSeries)
			t = t.Add(increment)
			x := 2.0 * math.Pi * (float64(j) / float64(pointsPerSeries))
			x += shift
			s.Data[j] = fluxcui.TimePoint{
				T: t,
				V: math.Sin(x),
			}
		}
		ss[i] = s
	}
	return ss
}

func (m Model) Timestamp() time.Time {
	t, err := time.Parse(time.RFC3339, "2018-06-26T13:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

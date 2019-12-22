package view

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/wolffcm/fluxcui/mock"
	"testing"
	"time"
)

func TestLineGraph_render(t *testing.T) {
	lg := newLineGraph()
	b := &bytes.Buffer{}
	if err := lg.render(mock.Model{}, 320, 96, b); err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%v", b.String())
}

func mustParseTime(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		panic(err)
	}

	return t
}

func TestLineGraph_chooseXTicks(t *testing.T) {
	tcs := []struct {
		name       string
		md         *Metadata
		wantXTicks []time.Time
	}{
		{
			name: "five minutes",
			md: &Metadata{
				minT: mustParseTime("2018-06-26T00:02:30Z"),
				maxT: mustParseTime("2018-06-26T00:07:30Z"),
				minV: 100,
				maxV: 1000,
			},
			wantXTicks: []time.Time{
				mustParseTime("2018-06-26T00:03:00Z"),
				mustParseTime("2018-06-26T00:04:00Z"),
				mustParseTime("2018-06-26T00:05:00Z"),
				mustParseTime("2018-06-26T00:06:00Z"),
				mustParseTime("2018-06-26T00:07:00Z"),
			},
		},
		{
			name: "one minute",
			md: &Metadata{
				minT: mustParseTime("2018-06-26T00:00:00Z"),
				maxT: mustParseTime("2018-06-26T00:01:00Z"),
				minV: 100,
				maxV: 1000,
			},
			wantXTicks: []time.Time{
				mustParseTime("2018-06-26T00:00:15Z"),
				mustParseTime("2018-06-26T00:00:30Z"),
				mustParseTime("2018-06-26T00:00:45Z"),
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			xts := chooseXTicks(tc.md)
			if diff := cmp.Diff(tc.wantXTicks, xts); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestLineGraph_chooseYTicks(t *testing.T) {
	tcs := []struct {
		name       string
		md         *Metadata
		wantYTicks []float64
	}{
		{
			name: "100 to 1000",
			md: &Metadata{
				minT: mustParseTime("2018-06-26T00:02:30Z"),
				maxT: mustParseTime("2018-06-26T00:07:30Z"),
				minV: 100,
				maxV: 1000,
			},
			wantYTicks: []float64{
				200, 300, 400, 500, 600, 700, 800, 900,
			},
		},
		{
			name: "-1 to 1",
			md: &Metadata{
				minT: mustParseTime("2018-06-26T00:02:30Z"),
				maxT: mustParseTime("2018-06-26T00:07:30Z"),
				minV: -1,
				maxV: 1,
			},
			wantYTicks: []float64{
				0,
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			xts := chooseYTicks(tc.md)
			if diff := cmp.Diff(tc.wantYTicks, xts); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

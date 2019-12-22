package view

import (
	"github.com/aybabtme/rgbterm"
	"testing"
)

func TestRGBTerm(t *testing.T) {
	numColors := 17
	for i := 0; i < numColors; i++ {
		f := float64(i) / float64(numColors-1)
		c := nineteenEightyFour.GetInterpolatedColorFor(f)
		r, g, b := c.RGB255()
		t.Log(rgbterm.FgString("hello world", r, g, b))
	}
}

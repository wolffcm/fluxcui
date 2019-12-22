package view

import (
	"github.com/aybabtme/rgbterm"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	grey = -1
)

func getPalette(numColors int) map[int]func(string) string {
	keypoints := GradientTable{
		{MustParseHex("#31c0f6"), 0.0},
		{MustParseHex("#ff7e27"), 0.5},
		{MustParseHex("#a500a5"), 1.0},
	}

	p := make(map[int]func(string) string)
	for i := 0; i < numColors; i++ {
		f := float64(i) / float64(numColors-1)
		c := keypoints.GetInterpolatedColorFor(f)
		r, g, b := c.RGB255()
		p[i] = func(s string) string {
			return rgbterm.FgString(s, r, g, b)
		}
	}
	p[grey] = func(s string) string {
		return rgbterm.FgString(s, 64, 64, 64)
	}
	return p
}

var nineteenEightyFour = GradientTable{
	{MustParseHex("#31c0f6"), 0.0},
	{MustParseHex("#ff7e27"), 0.5},
	{MustParseHex("#a500a5"), 1.0},
}

// Taken from github.com/lucasb-eyer/go-colorful gradient demo:

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

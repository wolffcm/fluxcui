package controller

import (
	"github.com/wolffcm/fluxcui"
	"github.com/wolffcm/fluxcui/model"
	"github.com/wolffcm/fluxcui/view"
)

type Config struct {
	Addr               string
	InsecureSkipVerify bool
	Token              string
}

type controller struct {
	cfg *Config

	m fluxcui.Model
	v fluxcui.View
}

func New(cfg *Config) (fluxcui.Controller, error) {
	c := &controller{
		cfg: cfg,
	}
	m, err := model.NewFluxModel(&model.Config{
		Addr:               cfg.Addr,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Token:              cfg.Token,
	})
	if err != nil {
		return nil, err
	}
	c.m = m

	vcfg := &view.Config{
		EditorText: `import "generate"
import "math"

sinWithShift = (v) => generate.from(
    start: 2018-06-26T00:00:00Z,
    stop: 2018-06-26T00:01:00Z,
    count: 256,
    fn: (n) => n
)
  |> map(fn: (r) => ({r with s: v, _value: math.sin(x: float(v: r._value) / 25.0 + v)}))
  |> group(columns: ["s"])

s0 = sinWithShift(v: 0.0 * 2.0 * math.pi / 5.0)
s1 = sinWithShift(v: 1.0 * 2.0 * math.pi / 5.0)
s2 = sinWithShift(v: 2.0 * 2.0 * math.pi / 5.0)
s3 = sinWithShift(v: 3.0 * 2.0 * math.pi / 5.0)
s4 = sinWithShift(v: 4.0 * 2.0 * math.pi / 5.0)
union(tables: [s0, s1, s2, s3, s4])
`,
	}
	v := view.NewView(vcfg, m, c)
	c.v = v
	return c, nil
}

func (c *controller) Run() error {
	return c.v.Run()
}

func (c *controller) Query(fluxSrc string) error {
	return c.m.Query(fluxSrc)
}

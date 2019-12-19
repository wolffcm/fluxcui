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
		EditorText: `from(bucket: "my-bucket")
  |> range(start: -5m)
  |> filter(fn: (r) => r._measurement == "diskio")
  |> filter(fn: (r) => r._field == "read_bytes" or r._field == "write_bytes")
  |> aggregateWindow(every: 10s, fn: last, createEmpty: false)
  |> derivative(unit: 1s, nonNegative: false)
  |> yield(name: "derivative")
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

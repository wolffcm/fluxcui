package model

import (
	"context"
	"errors"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/lang"
	"github.com/influxdata/flux/repl"
	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/http"
	"github.com/influxdata/influxdb/query"
	"github.com/wolffcm/fluxcui"
	"time"
)

type fluxData struct {
	querier repl.Querier
	deps flux.Dependencies

	ts time.Time
	s []fluxcui.Series
}

func NewFluxModel() (fluxcui.Model, error) {
	qs := &http.FluxQueryService{
		Addr:               "http://localhost:9999",
		Token:              "vjtk_WwsUAl8cNcgzPm7XWp8rZVtDiwShcqb3o5B7Zt-a5a358JY9eTo78-Y5B9RW2QfnSnGuHJp6pyE7knM6g==",
		InsecureSkipVerify: false,
	}
	orgID, err := influxdb.IDFromString("fbe7cf21e65601a1")
	if err != nil {
		return nil, err
	}
	q := &query.REPLQuerier{
		OrganizationID: *orgID,
		QueryService:   qs,
	}
	return &fluxData{
		querier: q,
		deps: flux.NewDefaultDependencies(),
	}, nil
}

func (f *fluxData) Timestamp() time.Time {
	return f.ts
}

func (f *fluxData) Query(fluxSrc string) error {
	ast, err := flux.Parse(fluxSrc)
	if err != nil {
		return err
	}

	c := lang.ASTCompiler{
		AST: ast,
		Now: time.Now(),
	}
	ri, err := f.querier.Query(context.TODO(), f.deps, c)
	if err != nil {
		return err
	}
	defer ri.Release()

	f.s = f.s[0:0]
	for ri.More() {
		r := ri.Next()
		ti := r.Tables()
		if err := ti.Do(func(t flux.Table) error {
			s := fluxcui.Series{
				Tags: make(map[string]string),
				Data: nil,
			}
			ti := execute.ColIdx("_time", t.Cols())
			if ti < 0 {
				return errors.New("missing _time column")
			}
			vi := execute.ColIdx("_value", t.Cols())
			if ti < 0 {
				return errors.New("missing _value column")
			}

			// TODO(cwolff): set up group key
			if err := t.Do(func(cr flux.ColReader) error {
				ts := cr.Times(ti)
				vs := cr.Floats(vi)
				for i := 0; i < cr.Len(); i++ {
					if ts.IsValid(i) && vs.IsValid(i) {
						timestamp := time.Unix(0, ts.Int64Values()[i])
						s.Data = append(s.Data, fluxcui.TimePoint{
							T: timestamp,
							V: vs.Float64Values()[i],
						})
					}
				}
				return nil
			}); err != nil {
				return err
			}
			f.s = append(f.s, s)
			return nil
		}); err != nil {
			return err
		}
	}
	if err := ri.Err(); err != nil {
		return err
	}
	f.ts = time.Now()
	return nil
}

func (f *fluxData) Series() []fluxcui.Series {
	return f.s
}




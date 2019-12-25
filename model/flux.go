package model

import (
	"context"
	"errors"
	"github.com/influxdata/flux"
	_ "github.com/influxdata/flux/builtin"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/lang"
	"github.com/influxdata/flux/memory"
	"github.com/influxdata/flux/repl"
	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/http"
	"github.com/influxdata/influxdb/query"
	"github.com/wolffcm/fluxcui"
	"time"
)

type Config struct {
	Addr               string
	InsecureSkipVerify bool
	Token              string
	OrgID              string
	Vanilla            bool
}

type fluxData struct {
	cfg     *Config
	querier repl.Querier
	deps    flux.Dependencies

	ts time.Time
	s  []fluxcui.Series
}

type vanillaQuerier struct{}

func (vanillaQuerier) Query(ctx context.Context, deps flux.Dependencies, c flux.Compiler) (flux.ResultIterator, error) {
	program, err := c.Compile(ctx)
	if err != nil {
		return nil, err
	}
	ctx = deps.Inject(ctx)
	alloc := &memory.Allocator{}
	qry, err := program.Start(ctx, alloc)
	if err != nil {
		return nil, err
	}
	return flux.NewResultIteratorFromQuery(qry), nil
}

func NewFluxModel(cfg *Config) (fluxcui.Model, error) {
	var q repl.Querier
	if !cfg.Vanilla {
		qs := &http.FluxQueryService{
			Addr:               cfg.Addr,
			Token:              cfg.Token,
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		}
		orgID, err := influxdb.IDFromString(cfg.OrgID)
		if err != nil {
			return nil, err
		}
		q = &query.REPLQuerier{
			OrganizationID: *orgID,
			QueryService:   qs,
		}
	} else {
		q = &vanillaQuerier{}
	}
	return &fluxData{
		querier: q,
		deps:    flux.NewDefaultDependencies(),
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
				ti = execute.ColIdx("_stop", t.Cols())
				if ti < 0 {
					return errors.New("missing _time or _stop column")
				}
			}
			vi := execute.ColIdx("_value", t.Cols())
			if ti < 0 {
				return errors.New("missing _value column")
			}
			vcm := t.Cols()[vi]
			vct := vcm.Type

			// TODO(cwolff): set up group key
			if err := t.Do(func(cr flux.ColReader) error {
				ts := cr.Times(ti)
				for i := 0; i < cr.Len(); i++ {
					var v float64
					var isValid bool
					switch vct {
					case flux.TFloat:
						v = cr.Floats(vi).Value(i)
						isValid = cr.Floats(vi).IsValid(i)
					case flux.TInt:
						v = float64(cr.Ints(vi).Value(i))
						isValid = cr.Ints(vi).IsValid(i)
					case flux.TUInt:
						v = float64(cr.UInts(vi).Value(i))
						isValid = cr.UInts(vi).IsValid(i)
					default:
					}
					if ts.IsValid(i) && isValid {
						timestamp := time.Unix(0, ts.Int64Values()[i])
						s.Data = append(s.Data, fluxcui.TimePoint{
							T: timestamp,
							V: v,
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

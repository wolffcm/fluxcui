package fluxcui

import (
	"time"
)

type TimePoint struct {
	T time.Time
	V float64
}

type Series struct {
	Tags map[string]string
	Data []TimePoint
}

type Model interface {
	Query(fluxSrc string) error
	Series() []Series
}

type View interface {
	Run() error
}

type Controller interface {
	Query(fluxSrc string) error
	Run() error
}
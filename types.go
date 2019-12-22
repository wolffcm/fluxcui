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

func IsEmpty(ss []Series) bool {
	for _, s := range ss {
		for _, _ = range s.Data {
			return false
		}
	}
	return true
}

type Model interface {
	Query(fluxSrc string) error
	Series() []Series
	Timestamp() time.Time
}

type View interface {
	Run() error
}

type Controller interface {
	Query(fluxSrc string) error
	Run() error
}

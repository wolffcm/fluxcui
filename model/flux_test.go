package model

import (
	"testing"
)

func TestFluxData_QueryVanilla(t *testing.T) {
	c := &Config{Vanilla: true}
	m, err := NewFluxModel(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Query(`
import "generate"
generate.from(
    start: 2018-06-26T00:00:00Z,
    stop: 2018-06-26T00:01:00Z,
    count: 16,
    fn: (n) => n
)
`); err != nil {
		t.Fatal(err)
	}
	ss := m.Series()
	if want, got := 1, len(ss); want != got {
		t.Fatalf("expected %v series, got %v", want, got)
	}
	if want, got := 16, len(ss[0].Data); want != got {
		t.Fatalf("expected %v points, got %v", want, got)
	}
}

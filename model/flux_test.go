package model

import (
	"testing"
)

func TestFluxData_Query(t *testing.T) {
	m, err := NewFluxModel()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Query(`
from(bucket: "my-bucket")
  |> range(start: -5m)
  |> filter(fn: (r) => r._measurement == "diskio")
  |> filter(fn: (r) => r._field == "read_bytes" or r._field == "write_bytes")
  |> aggregateWindow(every: 15s, fn: last, createEmpty: false)
  |> derivative(unit: 1s, nonNegative: false)
  |> yield(name: "derivative")
`); err != nil {
		t.Fatal(err)
	}
}

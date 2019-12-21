package view

import (
	"bytes"
	"github.com/wolffcm/fluxcui/mock"
	"testing"
)

func TestLineGraph_render(t *testing.T) {
	lg := newLineGraph()
	b := &bytes.Buffer{}
	if err := lg.render(mock.Model{}, 160, 48, b); err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%v", b.String())
}

package mock

import (
	"testing"
)

func TestModel_Series(t *testing.T) {
	m := Model{}
	ss := m.Series()
	s := ss[0]
	for _, pt := range s.Data {
		t.Logf("%#v", pt)
	}
}

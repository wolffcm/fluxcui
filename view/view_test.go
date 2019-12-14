package view

import (
	"fmt"
	"github.com/wolffcm/fluxcui/model"
	"testing"
	"time"
)

func TestScaleToCanvas(t *testing.T) {
	m := model.NewModel()
	tps := m.Series()[0].Data

	// Supposed a 25x25 view, with 25*2 x 25 * 4 drawille canvas
	cx, cy := 50.0, 100.0
	scaled := scaleToCanvas(tps, cx, cy)
	for i := 0; i < len(tps); i++ {
		fmt.Printf("(%v, %1.05f) -> (%0.5f, %0.5f)\n", tps[i].T.Format(time.RFC3339), tps[i].V, scaled[i].X, scaled[i].Y)
	}
}
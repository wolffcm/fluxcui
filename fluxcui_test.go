package fluxcui_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/wolffcm/fluxcui"
)

func TestScaleToCanvas(t *testing.T) {
	tps := fluxcui.GenData()

	// Supposed a 25x25 view, with 25*2 x 25 * 4 drawille canvas
	cx, cy := 50.0, 100.0
	scaled := fluxcui.ScaleToCanvas(tps, cx, cy)
	for i := 0; i < len(tps); i++ {
		fmt.Printf("(%v, %1.05f) -> (%0.5f, %0.5f)\n", tps[i].T.Format(time.RFC3339), tps[i].V, scaled[i].X, scaled[i].Y)
	}
}
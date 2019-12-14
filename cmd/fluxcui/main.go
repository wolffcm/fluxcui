package main

import (
	"fmt"
	"github.com/wolffcm/fluxcui"
	"log"
	"os"

	"github.com/exrook/drawille-go"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "fluxcui",
	Short: "Launch a Flux Console User Interface",
	Run: run,
}

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error executing command: %v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 1, 1, maxX - 2, maxY - 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := writeView(v); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func writeView(v *gocui.View) error {
	points := fluxcui.GenData()

	vxd, vyd := v.Size()
	cxd, cyd := vxd * 2, vyd * 4
	cps := fluxcui.ScaleToCanvas(points, float64(cxd), float64(cyd))
	c := drawille.NewCanvas()
	oldPt := cps[0]
	for _, pt := range cps[1:] {
		c.DrawLine(oldPt.X, oldPt.Y, pt.X, pt.Y)
		oldPt = pt
	}

	if _, err := v.Write([]byte(c.String())); err != nil {
		return err
	}

	//// inscribe a box in the view.
	//vxd, vyd := v.Size()
	//cxd, cyd := vxd * 2, vyd * 4


	return nil
}
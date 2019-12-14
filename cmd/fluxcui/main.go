package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wolffcm/fluxcui/model"
	"github.com/wolffcm/fluxcui/view"
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
	m := model.NewModel()
	v := view.NewView(m)
	if err := v.Run(); err != nil {
		panic(err)
	}
}

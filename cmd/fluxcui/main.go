package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wolffcm/fluxcui/controller"
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

func run(cmd *cobra.Command, _ []string) {
	c, err := controller.New(&controller.Config{})
	if err != nil {
		cmd.PrintErr(err)
	}
	if err := c.Run(); err != nil {
		cmd.PrintErr(err)
	}
}

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
	"github.com/wolffcm/fluxcui/controller"
)

var cmd = &cobra.Command{
	Use:   "fluxcui",
	Short: "Launch a Flux Console User Interface",
	Run:   run,
}

var (
	flagToken      string
	flagHost       string
	flagSkipVerify bool
)

func init() {
	viper.SetEnvPrefix("INFLUX")

	cmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "API token to be used throughout client calls")
	viper.BindEnv("TOKEN")
	if h := viper.GetString("TOKEN"); h != "" {
		flagToken = h
	}

	cmd.PersistentFlags().StringVar(&flagHost, "host", "http://localhost:9999", "HTTP address of InfluxDB")
	viper.BindEnv("HOST")
	if h := viper.GetString("HOST"); h != "" {
		flagHost = h
	}

	cmd.PersistentFlags().BoolVar(&flagSkipVerify, "skip-verify", false, "SkipVerify controls whether a client verifies the server's certificate chain and host name.")
}

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error executing command: %v", err)
	}
}

func run(cmd *cobra.Command, _ []string) {
	c, err := controller.New(&controller.Config{
		Addr:               flagHost,
		InsecureSkipVerify: flagSkipVerify,
		Token:              flagToken,
	})
	if err != nil {
		cmd.PrintErr(err)
	}
	if err := c.Run(); err != nil {
		cmd.PrintErr(err)
	}
}

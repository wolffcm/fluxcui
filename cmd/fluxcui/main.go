package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	flagOrgID      string
	flagVanilla    bool
	flagVerbose    bool
)

func init() {
	viper.SetEnvPrefix("INFLUX")

	cmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "API token to be used throughout client calls")
	if err := viper.BindEnv("TOKEN"); err != nil {
		panic(err)
	}
	if h := viper.GetString("TOKEN"); h != "" {
		flagToken = h
	}

	cmd.PersistentFlags().StringVarP(&flagOrgID, "org-id", "o", "", "organization ID")
	if err := viper.BindEnv("ORG_ID"); err != nil {
		panic(err)
	}
	if h := viper.GetString("ORG_ID"); h != "" {
		flagOrgID = h
	}

	cmd.PersistentFlags().StringVar(&flagHost, "host", "http://localhost:9999", "HTTP address of InfluxDB")
	if err := viper.BindEnv("HOST"); err != nil {
		panic(err)
	}
	if h := viper.GetString("HOST"); h != "" {
		flagHost = h
	}

	cmd.PersistentFlags().BoolVar(&flagSkipVerify, "skip-verify", false, "whether a client verifies the server's certificate chain and host name")

	cmd.PersistentFlags().BoolVar(&flagVanilla, "vanilla", false, `use "vanilla" Flux; don't connect to InfluxDB`)
	cmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, `log verbose output`)
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
		OrgID:              flagOrgID,
		Vanilla:            flagVanilla,
		Verbose:            flagVerbose,
	})
	if err != nil {
		cmd.PrintErr(err)
	}
	if err := c.Run(); err != nil {
		cmd.PrintErr(err)
	}
}

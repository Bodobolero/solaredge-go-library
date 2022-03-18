package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ulrichSchreiner/solaredge"
)

var (
	rootCmd = &cobra.Command{
		Use:   "solaredge",
		Short: "solaredge is a client for the solaredge webservice API",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	timezone string
)

func init() {
	t := time.Now()
	zone, _ := t.Zone()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("baseurl", solaredge.DEFAULT_URL, "The base URL for the webservices")
	rootCmd.PersistentFlags().StringVar(&timezone, "timezone", zone, "The timezone to use for timestamps")
	rootCmd.PersistentFlags().String("apikey", "", "Your API key")
	_ = viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	_ = viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
}

func initConfig() {
	viper.SetEnvPrefix("solaredge")
	viper.AutomaticEnv()
	if timezone != "" {
		_, err := time.LoadLocation(timezone)
		if err != nil {
			log.Fatalf("unknown timezone %q, please use name from the IANA time zone database: %v", timezone, err)
		}
		solaredge.SiteZone = timezone
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dumpAsJson(a any) string {
	d, _ := json.MarshalIndent(a, "", "  ")
	return string(d)
}

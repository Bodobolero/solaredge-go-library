package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ulrichSchreiner/solaredge"
)

var (
	siteCmd = &cobra.Command{
		Use:   "site",
		Short: "site related actions",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	detailsCmd = &cobra.Command{
		Use:   "details",
		Short: "query site details",
		Run: func(cmd *cobra.Command, args []string) {
			siteDetails()
		},
	}
	inventoryCmd = &cobra.Command{
		Use:   "inventory",
		Short: "query site inventory",
		Run: func(cmd *cobra.Command, args []string) {
			siteInventory()
		},
	}
	startTime   string
	endTime     string
	since       string
	storageData = &cobra.Command{
		Use:   "storagedata",
		Short: "query battery storage data",
		Run: func(cmd *cobra.Command, args []string) {
			start, end := getStartEnd()
			siteStorageData(start, end)
		},
	}
	powerDetails = &cobra.Command{
		Use:   "powerdetails",
		Short: "query power details",
		Run: func(cmd *cobra.Command, args []string) {
			start, end := getStartEnd()
			sitePowerDetails(start, end)
		},
	}
	energyDetails = &cobra.Command{
		Use:   "energydetails",
		Short: "query energy details",
		Run: func(cmd *cobra.Command, args []string) {
			start, end := getStartEnd()
			unit := solaredge.Quarter_Of_An_Hour
			if len(args) > 0 {
				unit = solaredge.TimeUnit(args[0])
			}
			siteEnergyDetails(unit, start, end)
		},
	}
)

func getStartEnd() (time.Time, time.Time) {
	dur, err := time.ParseDuration(since)
	if err != nil {
		log.Fatalf("cannot parse duration %q: %v", since, err)
	}
	start := time.Now().Add(dur * -1)
	end := time.Now()
	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			log.Fatalf("cannot parse %q as %q: %v", startTime, time.RFC3339, err)
		}
	}
	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			log.Fatalf("cannot parse %q as %q: %v", endTime, time.RFC3339, err)
		}
	}
	return start, end
}

func init() {
	siteCmd.PersistentFlags().String("siteid", "", "your site id to query")
	_ = viper.BindPFlag("siteid", siteCmd.PersistentFlags().Lookup("siteid"))
	storageData.PersistentFlags().StringVar(&startTime, "start", "", "the start time for the query or 1h in the past if empty, RFC3339")
	storageData.PersistentFlags().StringVar(&endTime, "end", "", "the end time for the query or 'now' if empty, RFC3339")
	storageData.PersistentFlags().StringVar(&since, "since", "1h", "the start of the query time range")
	powerDetails.PersistentFlags().StringVar(&startTime, "start", "", "the start time for the query or 1h in the past if empty, RFC3339")
	powerDetails.PersistentFlags().StringVar(&endTime, "end", "", "the end time for the query or 'now' if empty, RFC3339")
	powerDetails.PersistentFlags().StringVar(&since, "since", "1h", "the start of the query time range")
	energyDetails.PersistentFlags().StringVar(&startTime, "start", "", "the start time for the query or 1h in the past if empty, RFC3339")
	energyDetails.PersistentFlags().StringVar(&endTime, "end", "", "the end time for the query or 'now' if empty, RFC3339")
	energyDetails.PersistentFlags().StringVar(&since, "since", "1h", "the start of the query time range")
}

func siteClient() *solaredge.SiteClient {
	sic, err := solaredge.SiteFromIDs(viper.GetString("apikey"), viper.GetString("siteid"), solaredge.WithBaseURL(viper.GetString("baseurl")))
	if err != nil {
		log.Fatalf("cannot create client %v", err)
	}
	return sic
}

func siteDetails() {
	det, err := siteClient().Details()
	if err != nil {
		log.Fatalf("cannot query details: %v", err)
	}
	fmt.Printf("%s", dumpAsJson(det))
}

func siteInventory() {
	det, err := siteClient().Inventory()
	if err != nil {
		log.Fatalf("cannot query inventory: %v", err)
	}
	fmt.Printf("%s", dumpAsJson(det))
}

func siteStorageData(start, end time.Time) {
	det, err := siteClient().StorageData(start, end)
	if err != nil {
		log.Fatalf("cannot query storage data: %v", err)
	}
	fmt.Printf("%s", dumpAsJson(det))
}

func sitePowerDetails(start, end time.Time) {
	det, err := siteClient().PowerDetails(start, end)
	if err != nil {
		log.Fatalf("cannot query power details: %v", err)
	}
	fmt.Printf("%s", dumpAsJson(det))
}

func siteEnergyDetails(unit solaredge.TimeUnit, start, end time.Time) {
	det, err := siteClient().EnergyDetails(unit, start, end)
	if err != nil {
		log.Fatalf("cannot query energy details: %v", err)
	}
	fmt.Printf("%s", dumpAsJson(det))
}

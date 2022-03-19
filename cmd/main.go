package main

func main() {
	siteCmd.AddCommand(detailsCmd, inventoryCmd, storageData, powerDetails, energyDetails, powerflow)
	rootCmd.AddCommand(siteCmd)
	rootCmd.Execute()
}

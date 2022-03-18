package main

func main() {
	siteCmd.AddCommand(detailsCmd, inventoryCmd, storageData, powerDetails, energyDetails)
	rootCmd.AddCommand(siteCmd)
	rootCmd.Execute()
}

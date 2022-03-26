package main

func main() {
	siteCmd.AddCommand(detailsCmd, inventoryCmd, storageData, powerDetails, energyDetails, powerflow, overview)
	rootCmd.AddCommand(siteCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.Execute()
}

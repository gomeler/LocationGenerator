package cmd

import (
	"LocationGenerator/logging"

	"github.com/gomeler/LocationGenerator/generators"
	"github.com/spf13/cobra"
)

var regionNumberLocationsFlag int

func init() {
	regionCmd.Flags().IntVar(&regionNumberLocationsFlag, "locationCount", 0, "Number of locations to generate in the region.")
	rootCmd.AddCommand(regionCmd)
}

var regionCmd = &cobra.Command{
	Use:   "region",
	Short: "Generate a region full of locations",
	Run: func(cmd *cobra.Command, args []string) {
		//This is a clunky way of setting logging levels in different packages.
		if Verbose {
			logging.SetLevelDebug(log)
			generators.SetLevelDebug()
		} else {
			logging.SetLevelInfo(log)
			generators.SetLevelInfo()
		}
		region := generators.RegionEntry(regionNumberLocationsFlag)
		for _, location := range region.Locations {
			logLocation(location)
		}

	},
}

package cmd

import (
	"github.com/gomeler/LocationGenerator/generators"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(locationCmd)
}

var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Generate a location",
	Run: func(cmd *cobra.Command, args []string) {
		generators.LocationEntry()
	},
}

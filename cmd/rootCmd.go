package cmd

import (
	"os"

	"github.com/gomeler/LocationGenerator/logging"
	"github.com/spf13/cobra"
)

var log = logging.New()
var Verbose bool
var rootCmd = &cobra.Command{
	Use:   "ttgen",
	Short: "ttgen is a simple tabletop RPG location generator",
	Long:  `blerg`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

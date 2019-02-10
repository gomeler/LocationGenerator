package cmd

import (
	"fmt"
	"os"

	"github.com/gomeler/LocationGenerator/logging"
	"github.com/spf13/cobra"
)

var log = logging.New()

var rootCmd = &cobra.Command{
	Use:   "location",
	Short: "location is a simple tabletop RPG location generator",
	Long:  `blerg`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/gomeler/LocationGenerator/generators"
	"github.com/spf13/cobra"
)

var locationTypeFlag string

func init() {
	locationCmd.Flags().StringVar(&locationTypeFlag, "type", "random", "Type of location: [random, farm, hamlet, town, city]")
	rootCmd.AddCommand(locationCmd)
}

var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Generate a location",
	Run: func(cmd *cobra.Command, args []string) {
		//Boy, I hope I'm handling flags correctly with Cobra. Think I've set it up correctly with this local flag "type".
		var locationType int
		locationTypeFlag = strings.ToLower(locationTypeFlag)
		if locationTypeValidator(locationTypeFlag) {
			switch locationTypeFlag {
			case "random":
				locationType = rand.Intn(4)
			case "farm":
				locationType = 0
			case "hamlet":
				locationType = 1
			case "town":
				locationType = 2
			case "city":
				locationType = 3
			}
		} else {
			msg := fmt.Sprintf("Invalid input provided: %s", locationTypeFlag)
			log.Fatal(msg)
			panic(msg)
		}
		generators.LocationEntry(locationType)
	},
}

//Validate that the given flag is a valid choice.
func locationTypeValidator(flagInput string) bool {
	validOptions := []string{"random", "farm", "hamlet", "town", "city"}
	for _, val := range validOptions {
		if flagInput == val {
			return true
		}
	}
	return false
}

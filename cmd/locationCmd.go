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
		location := generators.LocationEntry(locationType)
		logLocation(location)
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

func logLocation(location []*generators.PopulatedBuilding) {
	//Going to need a prettyprint option for this.
	log.Info("Num buildings: ", len(location))
	var numPeople int
	var numNPCs int
	for _, building := range location {
		log.Info("Building: ", building.BaseBuilding.Name)
		if building.Owner.Name != "" {
			owner := fmt.Sprintf("%s %s %s %d", building.Owner.Name, building.Owner.Gender, building.Owner.Race, building.Owner.Age)
			log.Info(fmt.Sprintf("Owner: %s", owner))
			numPeople++
			numNPCs++
		}
		if building.Employees != nil {
			for _, employee := range building.Employees {
				e := fmt.Sprintf("%s %s %s %d", employee.Name, employee.Gender, employee.Race, employee.Age)
				log.Info(fmt.Sprintf("Subemployee: %s", e))
				numNPCs++
			}
			numPeople += len(building.Employees)
		}
		log.Info(fmt.Sprintf("Number of non-NPCs: %d", building.NonNPCEmployees))
		numPeople += building.NonNPCEmployees
	}
	log.Info(fmt.Sprintf("%d people in this location, %d NPCs", numPeople, numNPCs))
}

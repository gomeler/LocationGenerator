package cmd

import (
	"LocationGenerator/logging"
	"fmt"
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
		//This is a clunky way of setting logging levels in different packages.
		if Verbose {
			logging.SetLevelDebug(log)
			generators.SetLevelDebug()
		} else {
			logging.SetLevelInfo(log)
			generators.SetLevelInfo()
		}
		locationTypeFlag = strings.ToLower(locationTypeFlag)
		if locationTypeValidator(locationTypeFlag) {
			location := generators.LocationEntry(locationTypeFlag)
			logLocation(location)
		} else {
			msg := fmt.Sprintf("Invalid input provided: %s", locationTypeFlag)
			log.Fatal(msg)
			panic(msg)
		}
	},
}

//Validate that the given flag is a valid choice.
func locationTypeValidator(flagInput string) bool {
	//Programmatically retrieve the valid locations from TheLocations keys.
	validOptions := generators.GetLocationNames()
	validOptions = append(validOptions, "random")
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

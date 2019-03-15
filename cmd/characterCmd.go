package cmd

import (
	"LocationGenerator/logging"
	"fmt"
	"strings"

	"github.com/gomeler/LocationGenerator/generators"
	"github.com/spf13/cobra"
)

var characterRaceFlag string
var characterGenderFlag string

func init() {
	//Need to see if I can pull out the race names and feed them into the command info.
	characterCmd.Flags().StringVar(&characterRaceFlag, "race", "random", "Race for character, random by default.")
	characterCmd.Flags().StringVar(&characterGenderFlag, "gender", "random", "Gender for character, random by default.")
	rootCmd.AddCommand(characterCmd)
}

var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "Generate a character",
	Run: func(cmd *cobra.Command, args []string) {
		//This is a clunky way of setting logging levels in different packages.
		if Verbose {
			logging.SetLevelDebug(log)
			generators.SetLevelDebug()
		} else {
			logging.SetLevelInfo(log)
			generators.SetLevelInfo()
		}
		//Do better flag validation via looping.
		if raceFlagValidator(characterRaceFlag, generators.AssembleRaces()) {
			if genderFlagValidator(characterGenderFlag, generators.Genders) {
				//Awww yiss, the flags are valid.
				generators.CharacterEntry(characterRaceFlag, characterGenderFlag)
			} else {
				msg := fmt.Sprintf("Invalid input provided: %s", characterGenderFlag)
				log.Fatal(msg)
				panic(msg)
			}
		} else {
			msg := fmt.Sprintf("Invalid input provided: %s", characterRaceFlag)
			log.Fatal(msg)
			panic(msg)
		}
	},
}

func raceFlagValidator(flagInput string, races []generators.Race) bool {
	var validItems = make([]generators.WeightedItem, len(races))
	for idx, race := range races {
		validItems[idx] = race.WeightedItem
	}
	return genderFlagValidator(flagInput, validItems)
}

func genderFlagValidator(flagInput string, weightedItems []generators.WeightedItem) bool {
	validItems := generators.WeightedItemNames(weightedItems)
	validItems = append(validItems, "random")
	for _, val := range validItems {
		if strings.ToLower(flagInput) == strings.ToLower(val) {
			return true
		}
	}
	return false
}

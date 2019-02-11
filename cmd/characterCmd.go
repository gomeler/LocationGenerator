package cmd

import (
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
		//Do better flag validation via looping.
		if flagValidator(characterRaceFlag, generators.Races) {
			if flagValidator(characterGenderFlag, generators.Genders) {
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

func flagValidator(flagInput string, weightedItems []generators.WeightedItem) bool {
	validItems := generators.WeightedItemNames(weightedItems)
	validItems = append(validItems, "random")
	for _, val := range validItems {
		if strings.ToLower(flagInput) == strings.ToLower(val) {
			return true
		}
	}
	return false
}

package generators

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

//gender_race contains the base data used for generating a random gender and race. I think maybe there'll be various presets for race/gender in the future.

//The age related values in Race are used to configure a massaged normally distributed range with a minimum value. The idea is that this is a cheap/easy way to provide reasonable age generation without yielding things like an age of -42 years.
type Race struct {
	WeightedItem
	ageMean        float64
	ageSigma       float64
	ageSigmaFactor float64
}

//Should AssembleRaces be in an init value?
func AssembleRaces() []Race {
	var Races = make([]Race, 13)
	Races[0] = Race{WeightedItem: WeightedItem{Name: "Human", Weight: 125}, ageMean: 50, ageSigma: 25, ageSigmaFactor: 1.6}
	Races[1] = Race{WeightedItem: WeightedItem{Name: "Dragonborn", Weight: 5}, ageMean: 35, ageSigma: 25, ageSigmaFactor: 1.2}
	Races[2] = Race{WeightedItem: WeightedItem{Name: "Dwarf", Weight: 30}, ageMean: 175, ageSigma: 70, ageSigmaFactor: 2}
	Races[3] = Race{WeightedItem: WeightedItem{Name: "Elf", Weight: 10}, ageMean: 300, ageSigma: 175, ageSigmaFactor: 1.3}
	Races[4] = Race{WeightedItem: WeightedItem{Name: "Gnome", Weight: 12}, ageMean: 250, ageSigma: 130, ageSigmaFactor: 1.7}
	Races[5] = Race{WeightedItem: WeightedItem{Name: "Half Elf", Weight: 20}, ageMean: 90, ageSigma: 45, ageSigmaFactor: 1.6}
	Races[6] = Race{WeightedItem: WeightedItem{Name: "Half Orc", Weight: 7}, ageMean: 40, ageSigma: 25, ageSigmaFactor: 1.15}
	Races[7] = Race{WeightedItem: WeightedItem{Name: "Halfing", Weight: 15}, ageMean: 80, ageSigma: 60, ageSigmaFactor: 1.1}
	Races[8] = Race{WeightedItem: WeightedItem{Name: "Goblin", Weight: 3}, ageMean: 30, ageSigma: 20, ageSigmaFactor: 1.25}
	Races[9] = Race{WeightedItem: WeightedItem{Name: "Orc", Weight: 3}, ageMean: 25, ageSigma: 12, ageSigmaFactor: 1.2}
	Races[10] = Race{WeightedItem: WeightedItem{Name: "Tabaxi", Weight: 4}, ageMean: 50, ageSigma: 25, ageSigmaFactor: 1.6}
	Races[11] = Race{WeightedItem: WeightedItem{Name: "Warforged", Weight: 2}, ageMean: 600, ageSigma: 400, ageSigmaFactor: 1.35}
	Races[12] = Race{WeightedItem: WeightedItem{Name: "Half Giant", Weight: 1}, ageMean: 70, ageSigma: 40, ageSigmaFactor: 1.25}
	return Races
}

var Genders = []WeightedItem{
	WeightedItem{Name: "Male", Weight: 50},
	WeightedItem{Name: "Female", Weight: 50},
}

func raceArrayRandomWeightedSelect(races []Race) (int, error) {
	items := make([]WeightedItem, len(races))
	for idx, race := range races {
		items[idx] = race.WeightedItem
	}
	return RandomWeightedSelect(items)
}

//TODO: RandomRace and RandomGender are identical, merge with a new input.
//RandomRace performs a weighted select against the array of races provided above.
func RandomRace(characterRaceFlag string) (Race, error) {
	races := AssembleRaces()
	if characterRaceFlag == "random" {
		index, err := raceArrayRandomWeightedSelect(races)
		if err == nil {
			return races[index], err
		}
		return Race{}, err
	}
	return raceMatcher(characterRaceFlag, races)
}

//raceMatcher will look up a race that has been passed as a flag argument.
func raceMatcher(raceName string, races []Race) (Race, error) {
	for _, race := range races {
		if strings.ToLower(raceName) == strings.ToLower(race.Name) {
			return race, nil
		}
	}
	return Race{}, errors.New(fmt.Sprint("Failed to match race with : ", raceName))
}

//RandomGender performs a weighted select against the array of genders provided above.
func RandomGender(characterGenderFlag string) (string, error) {
	if characterGenderFlag == "random" {
		index, err := RandomWeightedSelect(Genders)
		if err == nil {
			return Genders[index].Name, err
		}
		return "", err
	} else {
		//Kick back the flag for now.
		//TODO: filter genders, return gender object?
		return characterGenderFlag, nil
	}
}

//WeightedItemNames exists primarily for input validation.
func WeightedItemNames(things []WeightedItem) []string {
	var names []string
	for _, thing := range things {
		names = append(names, thing.Name)
	}
	return names
}

func SemiNormalDistributionAgeGenerator(race Race) int {
	var sample = int(math.Round(rand.NormFloat64()*race.ageSigma + race.ageMean))
	if sample < int(math.Round(race.ageMean-(race.ageSigmaFactor*race.ageSigma))) {
		//generated too small of a sample. Retry.
		fmt.Println("Generated too small of a result: ", sample)
		sample = SemiNormalDistributionAgeGenerator(race)
	}
	return sample
}

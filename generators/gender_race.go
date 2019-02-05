package generators

//gender_race contains the base data used for generating a random gender and race. I think maybe there'll be various presets for race/gender in the future.

var races = []WeightedItem{
	WeightedItem{Name: "Human", Weight: 125},
	WeightedItem{Name: "Dragonborn", Weight: 5},
	WeightedItem{Name: "Dwarf", Weight: 30},
	WeightedItem{Name: "Elf", Weight: 10},
	WeightedItem{Name: "Gnome", Weight: 12},
	WeightedItem{Name: "Half Elf", Weight: 20},
	WeightedItem{Name: "Half Orc", Weight: 7},
	WeightedItem{Name: "Halfing", Weight: 15},
	WeightedItem{Name: "Goblin", Weight: 3},
	WeightedItem{Name: "Orc", Weight: 3},
	WeightedItem{Name: "Tabaxi", Weight: 4},
	WeightedItem{Name: "Warforged", Weight: 2},
	WeightedItem{Name: "Half Giant", Weight: 1},
}

var weightedRaces = SimpleWeightedItems{items: races}

var genders = []WeightedItem{
	WeightedItem{Name: "Male", Weight: 50},
	WeightedItem{Name: "Female", Weight: 50},
}

var weightedGenders = SimpleWeightedItems{items: genders}

//RandomRace performs a weighted select against the array of races provided above.
func RandomRace() (string, error) {
	index, err := weightedRaces.RandomWeightedSelect()
	if err == nil {
		return weightedRaces.items[index].Name, err
	}
	return "", err
}

//RandomGender performs a weighted select against the array of genders provided above.
func RandomGender() (string, error) {
	index, err := weightedGenders.RandomWeightedSelect()
	if err == nil {
		return weightedGenders.items[index].Name, err
	}
	return "", err
}

package generators

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

var genders = []WeightedItem{
	WeightedItem{Name: "Male", Weight: 50},
	WeightedItem{Name: "Female", Weight: 50},
}

func RandomRace() (string, error) {
	var weight int = TotalWeight(races)
	result, err := RandomWeightedSelect(races, weight)
	if err == nil {
		return result.Name, err
	}
	return "", err
}

func RandomGender() (string, error) {
	var weight int = TotalWeight(genders)
	result, err := RandomWeightedSelect(genders, weight)
	if err == nil {
		return result.Name, err
	}
	return "", err
}

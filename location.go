package main

import (
	"fmt"
	"math/rand"
	"time"

	"./generators"
)

func main() {
	/*adjective, err := generators.RandomAdjective()
	errorHandler(err)
	fmt.Println(adjective)*/
	buildingMap := generators.AssembleBuildings()
	rand.Seed(time.Now().UnixNano())
	var townSize int = rand.Intn(4)
	var numBuildings int
	var townWeight int
	switch townSize {
	case 0: //We're generating a Farm
		numBuildings = rand.Intn(generators.Farm.MaxBuildings-generators.Farm.MinBuildings) + generators.Farm.MinBuildings
		fmt.Println("Making a farm")
		townWeight = generators.Farm.Weight
	case 1: //We're generating a Hamlet
		numBuildings = rand.Intn(generators.Hamlet.MaxBuildings-generators.Hamlet.MinBuildings) + generators.Hamlet.MinBuildings
		fmt.Println("Making a hamlet")
		townWeight = generators.Hamlet.Weight
	case 2: //We're generating a Town
		numBuildings = rand.Intn(generators.Town.MaxBuildings-generators.Town.MinBuildings) + generators.Town.MinBuildings
		fmt.Println("Making a town")
		townWeight = generators.Town.Weight
	case 3: //We're generating a City
		numBuildings = rand.Intn(generators.City.MaxBuildings-generators.City.MinBuildings) + generators.City.MinBuildings
		fmt.Println("Making a city")
		townWeight = generators.City.Weight
	}
	fmt.Printf("We're generating: %d buildings\n", numBuildings)
	generateBuildings(numBuildings, townWeight, buildingMap)

}

func generateBuildings(numBuildings int, townWeight int, buildingMap map[string]*generators.WeightedBuildings) {
	rand.Seed(time.Now().UnixNano())

	//Keeping this logic outside of the main loop, we don't need to generate this array for each building.
	buildingTypes := make([]string, len(buildingMap))
	i := 0
	for k := range buildingMap {
		buildingTypes[i] = k
		i++
	}

	//This is our farm/town/city. Might make it an array of pointers to the original buildings.
	var buildings = make(map[string]int)
	var building generators.WeightedBuilding
	var generatedBuilding generators.WeightedBuilding
	//loop to generate each building.
	for i := 0; i < numBuildings; i++ {
		building = selectBuilding(townWeight, buildingMap, buildingTypes)
		fmt.Printf("Building selected: %s. ", building.Name)
		//check if building has child, so we can potentially evolve.
		if building.ChildBuilding != nil {
			generatedBuilding = building.EvolveBuilding(numBuildings, buildings)
		} else {
			//perform space checks.
			generatedBuilding = building
		}
		buildings[generatedBuilding.Name] += 1
		fmt.Printf("Added %s to the map.\n", generatedBuilding.Name)

	}
	fmt.Println(buildings)
}

//selectBuildingType arose from needing to move the building selection logic out of the primary loop in the chance that no valid building for the location existed in the selected WeightedBuilding array.
func selectBuilding(townWeight int,
	buildingMap map[string]*generators.WeightedBuildings,
	buildingTypes []string,
) generators.WeightedBuilding {

	var buildingTypeIdx int = rand.Intn(len(buildingTypes))
	var buildingTypeName string = buildingTypes[buildingTypeIdx]
	//Check to see if the selected type features a valid building for this townWeight
	if verifyBuildingTypeValid(townWeight, buildingMap[buildingTypeName]) {
		buildingIdx, err := buildingMap[buildingTypeName].RandomWeightedSelect()
		errorHandler(err)
		return buildingMap[buildingTypeName].Buildings[buildingIdx]
	}
	//Going to make a potentially catastrophic assumption here and assume that the provided townWeight will eventually yield a building.
	return selectBuilding(townWeight, buildingMap, buildingTypes)
}

//verifyBuildingTypeValid checks to see if the array of buildings to be selected from contains at least one option that will work with the given townWeight. This prevents a castle from spawning in a Farm location.
func verifyBuildingTypeValid(townWeight int, buildingTypeArray *generators.WeightedBuildings) bool {
	for i := 0; i < len(buildingTypeArray.Buildings); i++ {
		if buildingTypeArray.Buildings[i].MinCityWeight >= townWeight ||
			buildingTypeArray.Buildings[i].MinCityWeight == 0 {
			//Most WeightedBuildings will not have a MinCityWeight value set, meaning they're valid for all locations.
			return true
		}
	}
	return false
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

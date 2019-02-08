package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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
		building = selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, buildings)
		fmt.Printf("Building selected: %s. ", building.Name)
		//check if building has child, so we can potentially evolve.
		if building.ChildBuilding != nil {
			generatedBuilding = evolveBuilding(building, numBuildings, townWeight, buildings)
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
func selectBuilding(numBuildings int,
	townWeight int,
	buildingMap map[string]*generators.WeightedBuildings,
	buildingTypes []string,
	existingBuildings map[string]int,
) generators.WeightedBuilding {

	var buildingTypeIdx int = rand.Intn(len(buildingTypes))
	var buildingTypeName string = buildingTypes[buildingTypeIdx]
	var buildingName string
	//Check to see if the selected type features a valid building for this townWeight
	if verifyBuildingTypeValid(townWeight, buildingMap[buildingTypeName]) {
		buildingIdx, err := buildingMap[buildingTypeName].RandomWeightedSelect()
		errorHandler(err)
		buildingName = buildingMap[buildingTypeName].Buildings[buildingIdx].Name
		//Now check to see if the building fits within the townWeight
		if buildingMap[buildingTypeName].Buildings[buildingIdx].MinCityWeight <= townWeight ||
			buildingMap[buildingTypeName].Buildings[buildingIdx].MinCityWeight == 0 {
			fmt.Printf("%s with weight %d fits in townWeight %d. ",
				buildingName,
				buildingMap[buildingTypeName].Buildings[buildingIdx].MinCityWeight,
				townWeight)
			//Check to see if there is space in the town for the building.
			if verifyBuildingFits(buildingMap[buildingTypeName].Buildings[buildingIdx],
				existingBuildings, numBuildings) {
				fmt.Printf("%s fits within the town. ", buildingName)
				return buildingMap[buildingTypeName].Buildings[buildingIdx]
			}
			//Seriously just an excuse to test out strings.NewReplacer
			message := "{Building} does not fit within the town, insufficient space. There is a max of {MaxBuilding} buildings, and {Building} has a MaximumQuantity of {MaxQuantity} and a MaxPercentage of {MaxPercentage} with {NumBuildings} current {Building}. "
			r := strings.NewReplacer("{Building}", buildingName,
				"{MaxBuilding}", fmt.Sprint(numBuildings),
				"{MaxQuantity}", fmt.Sprint(buildingMap[buildingTypeName].Buildings[buildingIdx].MaxQuantity),
				"{MaxPercentage}", fmt.Sprint(buildingMap[buildingTypeName].Buildings[buildingIdx].MaximumPercentage),
				"{NumBuildings}", fmt.Sprint(existingBuildings[buildingName]))
			fmt.Printf(r.Replace(message))
		} else {
			fmt.Printf("%s does not fit in townWeight %d, rerolling. ", buildingName, townWeight)
			return selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, existingBuildings)
		}
	}
	//Going to make a potentially catastrophic assumption here and assume that the provided townWeight will eventually yield a building. In theory we'll just keep looping until we find a valid buildingType array that contains a building that fits within the townWeight value AND there is space for it in the existingBuilding map and the MaxQuantity/MaxPercentage checks.
	return selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, existingBuildings)
}

//verifyBuildingTypeValid checks to see if the array of buildings to be selected from contains at least one option that will work with the given townWeight. This prevents a castle from spawning in a Farm location.
func verifyBuildingTypeValid(townWeight int, buildingTypeArray *generators.WeightedBuildings) bool {
	for i := 0; i < len(buildingTypeArray.Buildings); i++ {
		if buildingTypeArray.Buildings[i].MinCityWeight <= townWeight ||
			buildingTypeArray.Buildings[i].MinCityWeight == 0 {
			//Most WeightedBuildings will not have a MinCityWeight value set, meaning they're valid for all locations.
			return true
		}
	}
	return false
}

//verifyBuildingFits checks to see if we haven't exceeded the MaximumPercentage and MaxQuantity values for the town and building combination.
func verifyBuildingFits(building generators.WeightedBuilding,
	existingBuildings map[string]int,
	maxBuildingCount int,
) bool {
	//Verify we won't exceed MaxQuantity. MaxQuantity is 0 by default, ignore those cases.
	if building.MaxQuantity < (existingBuildings[building.Name]+1) &&
		building.MaxQuantity != 0 {
		fmt.Printf("%s failed MaxQuantity check with MaxQuantity of %d and %d already available. ", building.Name, building.MaxQuantity, existingBuildings[building.Name])
		return false
	}
	//For exceedingly small locations a low MaximumPercentage basically means nothing spawns, so we'll loosen the limit and basically create a MinQuantity=1 by default.
	maxPercentageAppliedQuantity := float64(maxBuildingCount) * (float64(building.MaximumPercentage) / 100.0)
	if maxPercentageAppliedQuantity < float64(1) && maxPercentageAppliedQuantity > float64(0) {
		maxPercentageAppliedQuantity = 1.0
	}
	//Verify we won't exceed MaximumPercentage. MaximumPercentage is 0 by default, ignore those cases.
	if int(math.Round(maxPercentageAppliedQuantity)) < (existingBuildings[building.Name]+1) &&
		building.MaximumPercentage != 0 {
		fmt.Printf("%s failed MaxPercentage check with MaxPercentage of %d, an applied maximumQuantity of %d, and %d already available. ", building.Name, building.MaximumPercentage, int(math.Round(maxPercentageAppliedQuantity)), existingBuildings[building.Name])
		return false
	}
	return true
}

func evolveBuilding(wb generators.WeightedBuilding, maxBuildings int, townWeight int, buildingMap map[string]int) generators.WeightedBuilding {
	//Perform a roll to see if the childBuilding replaces the existing building. We return the child on success, otherwise we return the original value.
	var randomSelect int = rand.Intn(100) + 1
	if wb.ChildChance > randomSelect {
		//Success, evolution approved. Check if a grandchild exists and shortcircuit to that check.
		if wb.ChildBuilding.ChildBuilding == nil {
			//No grandchild, verify there is room for the child. Make sure the child also fits in the town, no castles in farms.
			fmt.Printf("Child evolution success: Evolved %s into %s. Checking if it fits. ", wb.Name, wb.ChildBuilding.Name)
			if verifyBuildingFits(*wb.ChildBuilding, buildingMap, maxBuildings) && (wb.ChildBuilding.MinCityWeight <= townWeight ||
				wb.ChildBuilding.MinCityWeight == 0) {
				fmt.Printf("Child fits. ")
				return *wb.ChildBuilding
			}
			fmt.Printf("Child does not fit, spawning parent. ")
			return wb
		} else {
			//grandchild exists, need to evolve it.
			//TODO: Need to handle the issue when there is no space for the child, there is space for the grandchild, and the child fails to evolve into the grandchild. This branch can yield more of a building than the limits dictate. Also should check if the grandchild can even fit in the town, can save some time if the evolution will auto-fail.
			fmt.Printf("Evolved %s to %s, now checking to see if we'll spawn a %s. ", wb.Name, wb.ChildBuilding.Name, wb.ChildBuilding.ChildBuilding.Name)
			return evolveBuilding(*wb.ChildBuilding, maxBuildings, townWeight, buildingMap)
		}
	}
	//Failed to evolve.
	fmt.Printf("Failed to evolve %s to %s. ", wb.Name, wb.ChildBuilding.Name)
	return wb
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

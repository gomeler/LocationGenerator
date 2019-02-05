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
	switch townSize {
	case 0: //We're generating a Farm
		numBuildings = rand.Intn(generators.Farm[1]-generators.Farm[0]) + generators.Farm[0]
		fmt.Println("Making a farm")
	case 1: //We're generating a Hamlet
		numBuildings = rand.Intn(generators.Hamlet[1]-generators.Hamlet[0]) + generators.Hamlet[0]
		fmt.Println("Making a hamlet")
	case 2: //We're generating a Town
		numBuildings = rand.Intn(generators.Town[1]-generators.Town[0]) + generators.Town[0]
		fmt.Println("Making a town")
	case 3: //We're generating a City
		numBuildings = rand.Intn(generators.City[1]-generators.City[0]) + generators.City[0]
		fmt.Println("Making a city")
	}
	fmt.Printf("We're generating: %d buildings\n", numBuildings)
	generateBuildings(numBuildings, buildingMap)

}

//map[string]int
func generateBuildings(numBuildings int, buildingMap map[string]*generators.WeightedBuildings) {
	//loop to generate each building.
	rand.Seed(time.Now().UnixNano())

	//Determine the building types. Kind of expecting all of the building types present in buildings.go to not necessarily be correct for various sized locations.
	buildingTypes := make([]string, len(buildingMap))
	i := 0
	for k := range buildingMap {
		buildingTypes[i] = k
		i++
	}

	var buildingType int
	var buildingTypeName string
	var buildingIdx int
	var buildingName string
	var err error
	var buildings = make(map[string]int)
	//var randomSelect int
	var generatedBuilding generators.WeightedBuilding
	for i := 1; i <= numBuildings; i++ {
		//randomly select a building type
		buildingType = rand.Intn(len(buildingTypes))
		buildingTypeName = buildingTypes[buildingType]
		fmt.Println(buildingTypeName)
		//fmt.Printf("Adding a %s\n", keys[buildingType])
		//buildingIdx = buildingMap[keys[buildingType]]
		buildingIdx, err = buildingMap[buildingTypeName].RandomWeightedSelect()
		errorHandler(err)
		buildingName = buildingMap[buildingTypeName].Buildings[buildingIdx].Name
		fmt.Printf("Building index: %d, building name: %s. ", buildingIdx, buildingName)
		//check if building has child

		if buildingMap[buildingTypeName].Buildings[buildingIdx].ChildBuilding != nil {
			//child evolution can happen
			generatedBuilding = buildingMap[buildingTypeName].Buildings[buildingIdx].EvolveBuilding(numBuildings, buildings)
		} else {
			//perform normal space checks
			//TODO: implement
			generatedBuilding = buildingMap[buildingTypeName].Buildings[buildingIdx]
		}
		buildings[generatedBuilding.Name] += 1
		fmt.Printf("Added %s to the map.\n", generatedBuilding.Name)

	}
	fmt.Println(buildings)
}

/*
func (wb *WeightedBuilding) evolveBuilding(maxBuildings int, buildingMap map[string]int) *WeightedBuilding {
	//We receive a pointer to a WeightedBuilding with a childBuilding value. We then perform a roll to see if the childBuilding replaces the existing building. We return the child on success, otherwise we return the original value.
	var randomSelect int = rand.Intn(100)
	var numExistingBuildings int //To be used
	if wb.childChance > randomSelect {
		//Success, evolution approved.
		if wb.childBuilding.childBuilding == nil {
			//No grandchild, verify there is room for the child
			if wb.childBuilding.maxQuantity {
				if buildingMap[wb.childBuilding.Name]+1 > wb.childBuilding.maxQuantity {
					//we'll have too many buildings, evolution denied.
					return wb
				}
			}
			if buildingMap[wb.childBuilding.Name]+1 > math.Round(maxBuildings*(wb.childBuilding.maximumPercentage/100)) {
				//There is insufficient space.
				return wb
			}
			//Passed all the checks, child can spawn.
			fmt.Printf("Evolved %s to %s. ", wb.Name, wb.childBuilding.Name)
			return wb.childBuilding
		} else {
			//grandchild exists, need to evolve it.
			fmt.Printf("Evolved %s to %s. ", wb.Name, wb.childBuilding.Name)
			return wb.childBuilding.evolveBuilding(maxBuildings, buildingMap)
		}
	}
	//Failed to evolve.
	return wb
}
*/

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

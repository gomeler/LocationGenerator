package generators

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/gomeler/LocationGenerator/logging"
)

var log = logging.New()

func LocationEntry() {
	/*adjective, err := generators.RandomAdjective()
	errorHandler(err)
	fmt.Println(adjective)*/
	buildingMap := AssembleBuildings()
	rand.Seed(time.Now().UnixNano())
	townSize := rand.Intn(4)
	var numBuildings int
	var townWeight int
	switch townSize {
	case 0: //We're generating a Farm
		numBuildings = rand.Intn(Farm.MaxBuildings-Farm.MinBuildings) + Farm.MinBuildings
		log.Info("Making a farm")
		townWeight = Farm.Weight
	case 1: //We're generating a Hamlet
		numBuildings = rand.Intn(Hamlet.MaxBuildings-Hamlet.MinBuildings) + Hamlet.MinBuildings
		log.Info("Making a hamlet")
		townWeight = Hamlet.Weight
	case 2: //We're generating a Town
		numBuildings = rand.Intn(Town.MaxBuildings-Town.MinBuildings) + Town.MinBuildings
		log.Info("Making a town")
		townWeight = Town.Weight
	case 3: //We're generating a City
		numBuildings = rand.Intn(City.MaxBuildings-City.MinBuildings) + City.MinBuildings
		log.Info("Making a city")
		townWeight = City.Weight
	}
	log.Info(fmt.Sprintf("We're generating: %d buildings.", numBuildings))
	generateBuildings(numBuildings, townWeight, buildingMap)

}

func generateBuildings(numBuildings int, townWeight int, buildingMap map[string]*WeightedBuildings) {
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
	var building WeightedBuilding
	var generatedBuilding WeightedBuilding
	//loop to generate each building.
	for i := 0; i < numBuildings; i++ {
		building = selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, buildings)
		log.Info(fmt.Sprintf("Building selected: %s.", building.Name))
		//check if building has child, so we can potentially evolve.
		generatedBuilding = evolveBuildingLoop(building, numBuildings, townWeight, buildings)
		buildings[generatedBuilding.Name]++
		log.Info(fmt.Sprintf("Added %s to the map.", generatedBuilding.Name))
	}
	log.Info(buildings)
}

//selectBuildingType arose from needing to move the building selection logic out of the primary loop in the chance that no valid building for the location existed in the selected WeightedBuilding array.
func selectBuilding(numBuildings int,
	townWeight int,
	buildingMap map[string]*WeightedBuildings,
	buildingTypes []string,
	existingBuildings map[string]int,
) WeightedBuilding {

	buildingTypeIdx := rand.Intn(len(buildingTypes))
	buildingTypeName := buildingTypes[buildingTypeIdx]
	var buildingName string
	//Check to see if the selected type features a valid building for this townWeight
	if verifyBuildingTypeValid(townWeight, buildingMap[buildingTypeName]) {
		buildingIdx, err := buildingMap[buildingTypeName].RandomWeightedSelect()
		errorHandler(err)
		buildingName = buildingMap[buildingTypeName].Buildings[buildingIdx].Name
		//Now check to see if the building fits within the townWeight
		if verifyTownWeight(buildingMap[buildingTypeName].Buildings[buildingIdx], townWeight) {
			log.Info(fmt.Sprintf("%s with weight %d fits in townWeight %d.",
				buildingName,
				buildingMap[buildingTypeName].Buildings[buildingIdx].MinCityWeight,
				townWeight))
			//Check to see if there is space in the town for the building.
			if verifyBuildingFits(buildingMap[buildingTypeName].Buildings[buildingIdx],
				existingBuildings, numBuildings) {
				log.Info(fmt.Sprintf("%s fits within the town.", buildingName))
				return buildingMap[buildingTypeName].Buildings[buildingIdx]
			}
			//Seriously just an excuse to test out strings.NewReplacer
			message := "{Building} does not fit within the town, insufficient space. There is a max of {MaxBuilding} buildings, and {Building} has a MaximumQuantity of {MaxQuantity} and a MaxPercentage of {MaxPercentage} with {NumBuildings} current {Building}."
			r := strings.NewReplacer("{Building}", buildingName,
				"{MaxBuilding}", fmt.Sprint(numBuildings),
				"{MaxQuantity}", fmt.Sprint(buildingMap[buildingTypeName].Buildings[buildingIdx].MaxQuantity),
				"{MaxPercentage}", fmt.Sprint(buildingMap[buildingTypeName].Buildings[buildingIdx].MaximumPercentage),
				"{NumBuildings}", fmt.Sprint(existingBuildings[buildingName]))
			log.Info(r.Replace(message))
		} else {
			log.Info(fmt.Sprintf("%s does not fit in townWeight %d, rerolling.", buildingName, townWeight))
			return selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, existingBuildings)
		}
	}
	//Going to make a potentially catastrophic assumption here and assume that the provided townWeight will eventually yield a building. In theory we'll just keep looping until we find a valid buildingType array that contains a building that fits within the townWeight value AND there is space for it in the existingBuilding map and the MaxQuantity/MaxPercentage checks.
	return selectBuilding(numBuildings, townWeight, buildingMap, buildingTypes, existingBuildings)
}

//verifyBuildingTypeValid checks to see if the array of buildings to be selected from contains at least one option that will work with the given townWeight.
func verifyBuildingTypeValid(townWeight int, buildingTypeArray *WeightedBuildings) bool {
	for i := 0; i < len(buildingTypeArray.Buildings); i++ {
		if verifyTownWeight(buildingTypeArray.Buildings[i], townWeight) {
			return true
		}
	}
	return false
}

//verifyTownWeight checks that the given building fits within the town. No castles in farms.
func verifyTownWeight(building WeightedBuilding, townWeight int) bool {
	if building.MinCityWeight <= townWeight || building.MinCityWeight == 0 {
		//Most WeightedBuildings will not have a MinCityWeight value set, meaning they're valid for all locations.
		return true
	}
	return false
}

//verifyBuildingFits checks to see if we haven't exceeded the MaximumPercentage and MaxQuantity values for the town and building combination.
func verifyBuildingFits(building WeightedBuilding,
	existingBuildings map[string]int,
	maxBuildingCount int,
) bool {
	//Verify we won't exceed MaxQuantity. MaxQuantity is 0 by default, ignore those cases.
	if building.MaxQuantity < (existingBuildings[building.Name]+1) &&
		building.MaxQuantity != 0 {
		log.Info(fmt.Sprintf("%s failed MaxQuantity check with MaxQuantity of %d and %d already available.", building.Name, building.MaxQuantity, existingBuildings[building.Name]))
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
		log.Info(fmt.Sprintf("%s failed MaxPercentage check with MaxPercentage of %d, an applied maximumQuantity of %d, and %d already available.", building.Name, building.MaximumPercentage, int(math.Round(maxPercentageAppliedQuantity)), existingBuildings[building.Name]))
		return false
	}
	return true
}

//evolveBuildingLoop is the second iteration of this evolution logic. I found it helpful to break apart the loop and upgrade logic. In some cases of chained evolutions, intermediate buildings don't necessarily need to fit, so we'll keep track of the most recent valid result and most recent evolution, and follow the evolutions to their conclusion.
func evolveBuildingLoop(building WeightedBuilding, maxBuildings int, townWeight int, buildingMap map[string]int) WeightedBuilding {
	//We'll store both the most recent valid placement, and the most recent placement. Non-zero chance the most recent placement won't fit, but may have a child building that fits.
	topBuilding := building  //Best building that fits.
	topEvolution := building //Best building regardless of fitment.
	var continueEvolution = true
	for continueEvolution && topEvolution.ChildBuilding != nil {
		//Check if we succeed in another evolution.
		continueEvolution, topEvolution = evolveBuildingCheck(topEvolution)
		//Check if the topEvolution fits in the town.
		if verifyBuildingFits(topEvolution, buildingMap, maxBuildings) &&
			verifyTownWeight(topEvolution, townWeight) {
			topBuilding = topEvolution
		}
	}
	if topBuilding != topEvolution {
		//This is strictly for my curiousity.
		log.Info(fmt.Sprintf("Turns out we rolled back %s to %s due to space.", topEvolution.Name, topBuilding.Name))
	}
	return topBuilding

}

//evolveBuildingCheck is strictly the evolution logic needed for one generation of evolution.
func evolveBuildingCheck(building WeightedBuilding) (bool, WeightedBuilding) {
	randomSelect := rand.Intn(100) + 1
	//Perform a roll to see if the childBuilding replaces the existing building. We return the child on success, otherwise we return the original value.
	if building.ChildChance > randomSelect {
		//Child evolution success, return the child.
		log.Info(fmt.Sprintf("Child evolution success: evolved %s into %s.", building.Name, building.ChildBuilding.Name))
		return true, *building.ChildBuilding
	}
	//Child evolution failure, we're done evolving things.
	log.Info(fmt.Sprintf("Child evolution failure for %s.", building.Name))
	return false, building

}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

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

type PopulatedBuilding struct {
	BaseBuilding    *WeightedBuilding
	Owner           NPC
	Employees       []NPC
	NonNPCEmployees int
	//Should probably account for things like wealth and inventory at some point.
}

func LocationEntry(townType int) {
	/*adjective, err := generators.RandomAdjective()
	errorHandler(err)
	fmt.Println(adjective)*/
	rand.Seed(time.Now().UnixNano())
	buildingMap := AssembleBuildings()
	var numBuildings int
	var townWeight int
	var location Location
	switch townType {
	case 0: //We're generating a Farm
		numBuildings = rand.Intn(Farm.MaxBuildings-Farm.MinBuildings) + Farm.MinBuildings
		log.Info("Making a farm")
		townWeight = Farm.Weight
		location = Farm
	case 1: //We're generating a Hamlet
		numBuildings = rand.Intn(Hamlet.MaxBuildings-Hamlet.MinBuildings) + Hamlet.MinBuildings
		log.Info("Making a hamlet")
		townWeight = Hamlet.Weight
		location = Hamlet
	case 2: //We're generating a Town
		numBuildings = rand.Intn(Town.MaxBuildings-Town.MinBuildings) + Town.MinBuildings
		log.Info("Making a town")
		townWeight = Town.Weight
		location = Town
	case 3: //We're generating a City
		numBuildings = rand.Intn(City.MaxBuildings-City.MinBuildings) + City.MinBuildings
		log.Info("Making a city")
		townWeight = City.Weight
		location = City
	}
	log.Info(fmt.Sprintf("We're generating: %d buildings.", numBuildings))
	town := generateBuildings(numBuildings, townWeight, buildingMap)
	log.Info(town)
	housing := make([]string, 4)
	housing[0] = "townhouse"
	housing[1] = "cottage"
	housing[2] = "shack"
	housing[3] = "hovel"
	for _, building := range town {
		//log.Info(*building.BaseBuilding)
		populateBuilding(&building, location.PeoplePerBuilding, location.NPCRatio)
		owner := fmt.Sprintf("%s %s %s %d", building.Owner.Name, building.Owner.Gender, building.Owner.Race, building.Owner.Age)
		if stringInSlice(building.BaseBuilding.Name, housing) {
			//pass
		} else {
			log.Info(fmt.Sprintf("Building: %s - Owner: %s - Num Employees: %d", building.BaseBuilding.Name, owner, building.NonNPCEmployees))
			if building.Employees != nil {
				for _, employee := range building.Employees {
					e := fmt.Sprintf("%s %s %s %d", employee.Name, employee.Gender, employee.Race, employee.Age)
					log.Info(fmt.Sprintf("Subemployee: %s", e))
				}
			}
		}
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func generateBuildings(numBuildings int, townWeight int, buildingMap []*WeightedBuildingCollection) []PopulatedBuilding {
	rand.Seed(time.Now().UnixNano())

	//buildings will be used for fitment checks. I think it'll be more efficient to track building quantities in a separate map vs having to iterate over the array of buildings and extract their type.
	var buildings = make(map[string]int)
	var discreteBuildings = make([]PopulatedBuilding, numBuildings)
	var building *WeightedBuilding
	var generatedBuilding *WeightedBuilding
	//loop to generate each building.
	for i := 0; i < numBuildings; i++ {
		building = selectBuilding(numBuildings, townWeight, buildingMap, buildings)
		log.Info(fmt.Sprintf("Building selected: %s.", building.Name))
		//check if building has child, so we can potentially evolve.
		generatedBuilding = evolveBuildingLoop(*building, numBuildings, townWeight, buildings)
		buildings[generatedBuilding.Name]++
		log.Info(fmt.Sprintf("Added %s to the map.", generatedBuilding.Name))
		discreteBuildings[i] = PopulatedBuilding{BaseBuilding: generatedBuilding}
	}
	log.Info(buildings)
	return discreteBuildings
}

//selectBuildingType arose from needing to move the building selection logic out of the primary loop in the chance that no valid building for the location existed in the selected WeightedBuilding array.
func selectBuilding(numBuildings int,
	townWeight int,
	buildingCollections []*WeightedBuildingCollection,
	existingBuildings map[string]int,
) *WeightedBuilding {

	buildingCollectionIdx, err := WeightedBuildingCollectionRandomWeightedSelect(buildingCollections)
	errorHandler(err)
	var candidateBuilding WeightedBuilding
	//Check to see if the selected buildingCollection features a valid building for this townWeight
	if verifyBuildingTypeValid(townWeight, buildingCollections[buildingCollectionIdx]) {
		buildingIdx, err := BuildingsRandomWeightedSelect(buildingCollections[buildingCollectionIdx].Buildings)
		errorHandler(err)
		candidateBuilding = buildingCollections[buildingCollectionIdx].Buildings[buildingIdx]
		//Now check to see if the selected building fits within the townWeight. We could avoid this check by feeding BuildingsRandomWeightedSelect with the townWeight so we don't select an invalid option.
		if verifyTownWeight(candidateBuilding, townWeight) {
			log.Info(fmt.Sprintf("%s with weight %d fits in townWeight %d.",
				candidateBuilding.Name,
				candidateBuilding.MinCityWeight,
				townWeight))
			//Check to see if there is space in the town for the building.
			if verifyBuildingFits(buildingCollections[buildingCollectionIdx].Buildings[buildingIdx],
				existingBuildings, numBuildings) {
				log.Info(fmt.Sprintf("%s fits within the town.", candidateBuilding.Name))
				return &buildingCollections[buildingCollectionIdx].Buildings[buildingIdx]
			}
			//Seriously just an excuse to test out strings.NewReplacer
			message := "{Building} does not fit within the town, insufficient space. There is a max of {MaxBuilding} buildings, and {Building} has a MaximumQuantity of {MaxQuantity} and a MaxPercentage of {MaxPercentage} with {NumBuildings} current {Building}."
			r := strings.NewReplacer("{Building}", candidateBuilding.Name,
				"{MaxBuilding}", fmt.Sprint(numBuildings),
				"{MaxQuantity}", fmt.Sprint(candidateBuilding.MaxQuantity),
				"{MaxPercentage}", fmt.Sprint(candidateBuilding.MaximumPercentage),
				"{NumBuildings}", fmt.Sprint(existingBuildings[candidateBuilding.Name]))
			log.Info(r.Replace(message))
		} else {
			log.Info(fmt.Sprintf("%s does not fit in townWeight %d, rerolling.", candidateBuilding.Name, townWeight))
			return selectBuilding(numBuildings, townWeight, buildingCollections, existingBuildings)
		}
	}
	//Going to make a potentially catastrophic assumption here and assume that the provided townWeight will eventually yield a building. In theory we'll just keep looping until we find a valid buildingType array that contains a building that fits within the townWeight value AND there is space for it in the existingBuilding map and the MaxQuantity/MaxPercentage checks.
	return selectBuilding(numBuildings, townWeight, buildingCollections, existingBuildings)
}

//verifyBuildingTypeValid checks to see if the array of buildings to be selected from contains at least one option that will work with the given townWeight.
func verifyBuildingTypeValid(townWeight int, buildingTypeArray *WeightedBuildingCollection) bool {
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
func evolveBuildingLoop(building WeightedBuilding, maxBuildings int, townWeight int, buildingMap map[string]int) *WeightedBuilding {
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
	return &topBuilding
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

func populateBuilding(building *PopulatedBuilding, peoplePerBuilding int, townNPCRatio float64) {
	if building.BaseBuilding.HasOwner {
		building.Owner = CharacterAttachment("random", "random")
	}
	//Determine number of employees
	var numEmployees int
	if building.BaseBuilding.MaxNumEmployees != 0 && building.BaseBuilding.MinNumEmployees != 0 {
		numEmployees = rand.Intn(building.BaseBuilding.MaxNumEmployees+1-building.BaseBuilding.MinNumEmployees) + building.BaseBuilding.MinNumEmployees
		numNPCs := int(math.Round(float64(numEmployees) * townNPCRatio))
		if numNPCs > 0 {
			for idx := 0; idx < numNPCs; idx++ {
				building.Employees = append(building.Employees, CharacterAttachment("random", "random"))
			}

		} else {
			//feature a random chance to spawn a singular NPC employee. Check to make sure we won't exceed min/max employee counts.
		}
		building.NonNPCEmployees = numEmployees - numNPCs
	}
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

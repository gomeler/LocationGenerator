package generators

import (
	"LocationGenerator/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gomeler/LocationGenerator/logging"
)

var log = logging.New()

func SetLevelDebug() {
	logging.SetLevelDebug(log)
}

func SetLevelInfo() {
	logging.SetLevelInfo(log)
}

type Location struct {
	WeightedItem
	MinBuildings      int
	MaxBuildings      int
	GrowthFactor      float64 //To be used with events to age location.
	NPCRatio          float64 //Ratio designed to reduce the number of pre-generated NPCs. To be adjusted in the future.
	PeoplePerBuilding int     //Settling on 8 people/building unless defined in the building type.
	TownTier          int
}

// Region is the largest unit we'll work with for now, and what represents the space that perhaps a kingdom or wilderness area would occupy. Regions will have master values for factions, existing wealth, population growth, and wealth growth that will be applied to underlying locations.
type Region struct {
	Name         string
	Locations    [][]*PopulatedBuilding
	Factions     string  //This is a stand-in for the regional faction system.
	GrowthFactor float64 //To be used with events to age regions.
}

//TheLocations consists of a map of available locations. On init we'll load a default value that can be configured by the CLI, or programmatically from other functions. Not sure if we should have these global values present, seems like a good way to cause a mistake.
var TheLocations = make(map[string]Location)

func init() {
	ReloadLocationFile("")
}

//ReloadLocationFile loads the default location config file when passed an empty string. Otherwise it loads the file it is pointed to and replaces TheLocations variable with the new content.
func ReloadLocationFile(filename string) {
	if filename == "" {
		//Empty string = load the project provided default.
		filename = config.LocationConfig
	}
	loadedFile, err := ioutil.ReadFile(filename)
	errorHandler(err)
	json.Unmarshal(loadedFile, &TheLocations)
}

type PopulatedBuilding struct {
	BaseBuilding    *WeightedBuilding
	Owner           NPC
	Employees       []NPC
	NonNPCEmployees int
	//Should probably account for things like wealth and inventory at some point.
}

// RegionEntry is a temporary stand-in for the development of the region system. I think a Region struct likely needs to be built for the return value, and as a means to more easily manage multiple regions.
func RegionEntry(numberLocations int) Region {
	var locations [][]*PopulatedBuilding = make([][]*PopulatedBuilding, numberLocations)
	for i := 0; i < numberLocations; i++ {
		locations[i] = LocationEntry("weightedrandom")
	}
	var region Region = Region{"Temporary name", locations, "temporary faction data", 0.0}
	return region
}

func LocationEntry(locationName string) []*PopulatedBuilding {
	/*adjective, err := generators.RandomAdjective()
	errorHandler(err)
	fmt.Println(adjective)*/
	rand.Seed(time.Now().UnixNano())
	var numBuildings int
	location, error := selectLocation(locationName)
	log.Info("Making a ", location.Name)
	errorHandler(error)
	numBuildings = rand.Intn(location.MaxBuildings-location.MinBuildings) + location.MinBuildings
	log.Info(fmt.Sprintf("We're generating: %d buildings.", numBuildings))
	spawnedLocation := generateBuildings(numBuildings, location.TownTier, TheBuildings)
	//Go about populating the buildings. For the most part housing won't have any NPCs, businesses will have an owner and 0,n employees. Buildings will also have a small chance at spawning a random NPC just for giggles if they don't generate one via the town NPC ratio check.
	for _, building := range spawnedLocation {
		populateBuilding(building, location.PeoplePerBuilding, location.NPCRatio)
	}
	return spawnedLocation
}

func GetLocationNames() []string {
	locationNames := make([]string, len(TheLocations))
	var idx int
	for key := range TheLocations {
		locationNames[idx] = key
		idx++
	}
	return locationNames
}

//SelectRandomLocation will pull a random location from TheLocations, not factoring any weighting.
func SelectRandomLocation() Location {
	var locationNames = GetLocationNames()
	var idx int = rand.Intn(len(locationNames))
	return TheLocations[locationNames[idx]]
}

func selectLocation(locationName string) (Location, error) {
	if locationName == "random" {
		return SelectRandomLocation(), nil
	} else if locationName == "weightedrandom" {
		//Perform a weighted select to retrieve a location name, and then fall through to the confirmation and return of the location below.
		selectedLocationName, error := LocationRandomWeightedSelect(TheLocations)
		errorHandler(error)
		locationName = selectedLocationName
	}
	if location, ok := TheLocations[locationName]; ok {
		return location, nil
	}
	return Location{}, fmt.Errorf("Failed to match location name %s", locationName)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func generateBuildings(numBuildings int, townTier int, buildingMap []*WeightedBuildingCollection) []*PopulatedBuilding {
	rand.Seed(time.Now().UnixNano())

	//buildings will be used for fitment checks. I think it'll be more efficient to track building quantities in a separate map vs having to iterate over the array of buildings and extract their type.
	var buildings = make(map[string]int)
	var discreteBuildings = make([]*PopulatedBuilding, numBuildings)
	var building *WeightedBuilding
	var generatedBuilding *WeightedBuilding
	//loop to generate each building.
	for i := 0; i < numBuildings; i++ {
		building = selectBuilding(numBuildings, townTier, buildingMap, buildings)
		log.Debug(fmt.Sprintf("Building selected: %s.", building.Name))
		//check if building has child, so we can potentially evolve.
		generatedBuilding = evolveBuildingLoop(*building, numBuildings, townTier, buildings)
		buildings[generatedBuilding.Name]++
		log.Debug(fmt.Sprintf("Added %s to the map.", generatedBuilding.Name))
		discreteBuildings[i] = &PopulatedBuilding{BaseBuilding: generatedBuilding}
	}
	log.Debug(buildings)
	return discreteBuildings
}

//selectBuildingType arose from needing to move the building selection logic out of the primary loop in the chance that no valid building for the location existed in the selected WeightedBuilding array.
func selectBuilding(numBuildings int,
	townTier int,
	buildingCollections []*WeightedBuildingCollection,
	existingBuildings map[string]int,
) *WeightedBuilding {

	buildingCollectionIdx, err := WeightedBuildingCollectionRandomWeightedSelect(buildingCollections)
	errorHandler(err)
	var candidateBuilding WeightedBuilding
	//Check to see if the selected buildingCollection features a valid building for this townTier
	if verifyBuildingTypeValid(townTier, buildingCollections[buildingCollectionIdx]) {
		buildingIdx, err := BuildingsRandomWeightedSelect(buildingCollections[buildingCollectionIdx].Buildings)
		errorHandler(err)
		candidateBuilding = buildingCollections[buildingCollectionIdx].Buildings[buildingIdx]
		//Now check to see if the selected building fits within the townTier. We could avoid this check by feeding BuildingsRandomWeightedSelect with the townTier so we don't select an invalid option.
		if verifytownTier(candidateBuilding, townTier) {
			log.Debug(fmt.Sprintf("%s with tier %d fits in townTier %d.",
				candidateBuilding.Name,
				candidateBuilding.MinCityWeight,
				townTier))
			//Check to see if there is space in the town for the building.
			if verifyBuildingFits(buildingCollections[buildingCollectionIdx].Buildings[buildingIdx],
				existingBuildings, numBuildings) {
				log.Debug(fmt.Sprintf("%s fits within the town.", candidateBuilding.Name))
				return &buildingCollections[buildingCollectionIdx].Buildings[buildingIdx]
			}
			//Seriously just an excuse to test out strings.NewReplacer
			message := "{Building} does not fit within the town, insufficient space. There is a max of {MaxBuilding} buildings, and {Building} has a MaximumQuantity of {MaxQuantity} and a MaxPercentage of {MaxPercentage} with {NumBuildings} current {Building}."
			r := strings.NewReplacer("{Building}", candidateBuilding.Name,
				"{MaxBuilding}", fmt.Sprint(numBuildings),
				"{MaxQuantity}", fmt.Sprint(candidateBuilding.MaxQuantity),
				"{MaxPercentage}", fmt.Sprint(candidateBuilding.MaximumPercentage),
				"{NumBuildings}", fmt.Sprint(existingBuildings[candidateBuilding.Name]))
			log.Debug(r.Replace(message))
		} else {
			log.Debug(fmt.Sprintf("%s does not fit in townTier %d, rerolling.", candidateBuilding.Name, townTier))
			return selectBuilding(numBuildings, townTier, buildingCollections, existingBuildings)
		}
	}
	//Going to make a potentially catastrophic assumption here and assume that the provided townTier will eventually yield a building. In theory we'll just keep looping until we find a valid buildingType array that contains a building that fits within the townTier value AND there is space for it in the existingBuilding map and the MaxQuantity/MaxPercentage checks.
	return selectBuilding(numBuildings, townTier, buildingCollections, existingBuildings)
}

//verifyBuildingTypeValid checks to see if the array of buildings to be selected from contains at least one option that will work with the given townTier.
func verifyBuildingTypeValid(townTier int, buildingTypeArray *WeightedBuildingCollection) bool {
	for i := 0; i < len(buildingTypeArray.Buildings); i++ {
		if verifytownTier(buildingTypeArray.Buildings[i], townTier) {
			return true
		}
	}
	return false
}

//verifytownTier checks that the given building fits within the town. No castles in farms.
func verifytownTier(building WeightedBuilding, townTier int) bool {
	if building.MinCityWeight <= townTier || building.MinCityWeight == 0 {
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
		log.Debug(fmt.Sprintf("%s failed MaxQuantity check with MaxQuantity of %d and %d already available.", building.Name, building.MaxQuantity, existingBuildings[building.Name]))
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
		log.Debug(fmt.Sprintf("%s failed MaxPercentage check with MaxPercentage of %d, an applied maximumQuantity of %d, and %d already available.", building.Name, building.MaximumPercentage, int(math.Round(maxPercentageAppliedQuantity)), existingBuildings[building.Name]))
		return false
	}
	return true
}

//evolveBuildingLoop is the second iteration of this evolution logic. I found it helpful to break apart the loop and upgrade logic. In some cases of chained evolutions, intermediate buildings don't necessarily need to fit, so we'll keep track of the most recent valid result and most recent evolution, and follow the evolutions to their conclusion.
func evolveBuildingLoop(building WeightedBuilding, maxBuildings int, townTier int, buildingMap map[string]int) *WeightedBuilding {
	//We'll store both the most recent valid placement, and the most recent placement. Non-zero chance the most recent placement won't fit, but may have a child building that fits.
	topBuilding := building  //Best building that fits.
	topEvolution := building //Best building regardless of fitment.
	var continueEvolution = true
	for continueEvolution && topEvolution.ChildBuilding != nil {
		//Check if we succeed in another evolution.
		continueEvolution, topEvolution = evolveBuildingCheck(topEvolution)
		//Check if the topEvolution fits in the town.
		if verifyBuildingFits(topEvolution, buildingMap, maxBuildings) &&
			verifytownTier(topEvolution, townTier) {
			topBuilding = topEvolution
		}
	}
	if topBuilding != topEvolution {
		//This is strictly for my curiousity.
		log.Debug(fmt.Sprintf("Turns out we rolled back %s to %s due to space.", topEvolution.Name, topBuilding.Name))
	}
	return &topBuilding
}

//evolveBuildingCheck is strictly the evolution logic needed for one generation of evolution.
func evolveBuildingCheck(building WeightedBuilding) (bool, WeightedBuilding) {
	randomSelect := rand.Intn(100) + 1
	//Perform a roll to see if the childBuilding replaces the existing building. We return the child on success, otherwise we return the original value.
	if building.ChildChance > randomSelect {
		//Child evolution success, return the child.
		log.Debug(fmt.Sprintf("Child evolution success: evolved %s into %s.", building.Name, building.ChildBuilding.Name))
		return true, *building.ChildBuilding
	}
	//Child evolution failure, we're done evolving things.
	log.Debug(fmt.Sprintf("Child evolution failure for %s.", building.Name))
	return false, building
}

//Breaking up building population generation for testing purposes.
func spawnOwnerForBuilding(building *PopulatedBuilding) {
	//This is a pretty simple task. Buildings have a bool flag that dictates if they have an owner.
	if building.BaseBuilding.HasOwner {
		building.Owner = CharacterAttachment("random", "random")
	}
}

func spawnEmployeesForBuilding(building *PopulatedBuilding, peoplePerBuilding int, townNPCRatio float64) {
	//Some buildings will have min/max employee ranges. Otherwise we'll use the town's default value.
	var numEmployees int
	if building.BaseBuilding.MaxNumEmployees != 0 && building.BaseBuilding.MinNumEmployees != 0 {
		numEmployees = rand.Intn(building.BaseBuilding.MaxNumEmployees+1-building.BaseBuilding.MinNumEmployees) + building.BaseBuilding.MinNumEmployees
	} else {
		numEmployees = rand.Intn(peoplePerBuilding + 1)
	}
	//Of the employees generated, some of them will be spawned as full NPCs.
	numNPCs := int(math.Round(float64(numEmployees) * townNPCRatio))
	//log.Info(fmt.Sprintf("%s will have %d employees, of which %d will be NPCs.", building.BaseBuilding.Name, numEmployees, numNPCs))
	building.NonNPCEmployees = numEmployees - numNPCs
	if numNPCs > 0 {
		for idx := 0; idx < numNPCs; idx++ {
			building.Employees = append(building.Employees, CharacterAttachment("random", "random"))
		}
	} else {
		//Feature a random chance to spawn a singular NPC employee. This exists as a second chance of sorts so towns don't end up simply with building owners and otherwise empty buildings. Check to make sure we won't exceed min/max employee counts.
		if rand.Intn(100) > 84 && numEmployees < building.BaseBuilding.MaxNumEmployees {
			building.Employees = append(building.Employees, CharacterAttachment("random", "random"))
		}
	}
}

func populateBuilding(building *PopulatedBuilding, peoplePerBuilding int, townNPCRatio float64) {
	//spawnOwnerForBuilding contains an internal check to determine if a building has an owner.
	spawnOwnerForBuilding(building)

	//spawnEmployeesForBuilding spawns NPC and non-NPC employees per the building's provided parameters. Otherwise it works against the default values provided by the town.
	spawnEmployeesForBuilding(building, peoplePerBuilding, townNPCRatio)
}

//Bunch of stuff needed to save/load Location structs to/from disk as JSON encoded objects.
type JSONLocation struct {
	Name              string  `json:"name"`
	Weight            int     `json:"weight"`
	MinBuildings      int     `json:"minBuildings"`
	MaxBuildings      int     `json:"maxBuildings"`
	GrowthFactor      float64 `json:"growthFactor"`
	NPCRatio          float64 `json:"npcRatio"`
	PeoplePerBuilding int     `json:"peoplePerBuilding"`
	TownTier          int     `json:"townTier"`
}

func (location Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONLocation(location))
}

func (location *Location) UnmarshalJSON(data []byte) error {
	var jsonLocation JSONLocation
	if err := json.Unmarshal(data, &jsonLocation); err != nil {
		return err
	}
	*location = jsonLocation.TranslateLocation()
	return nil
}

func NewJSONLocation(location Location) JSONLocation {
	return JSONLocation{location.WeightedItem.Name, location.WeightedItem.Weight, location.MinBuildings, location.MaxBuildings, location.GrowthFactor, location.NPCRatio, location.PeoplePerBuilding, location.TownTier}
}

func (jsonLocation JSONLocation) TranslateLocation() Location {
	location := Location{}
	location.WeightedItem.Name = jsonLocation.Name
	location.WeightedItem.Weight = jsonLocation.Weight
	location.MinBuildings = jsonLocation.MinBuildings
	location.MaxBuildings = jsonLocation.MaxBuildings
	location.GrowthFactor = jsonLocation.GrowthFactor
	location.NPCRatio = jsonLocation.NPCRatio
	location.PeoplePerBuilding = jsonLocation.PeoplePerBuilding
	location.TownTier = jsonLocation.TownTier
	return location
}

func tempWriteLocation(filename string) {
	file, err := os.Create(filename)
	errorHandler(err)
	defer file.Close()
	var jsonEncodedBuildings []byte
	jsonEncodedBuildings, err = json.Marshal(TheLocations)
	errorHandler(err)
	_, err = file.Write(jsonEncodedBuildings)
	errorHandler(err)
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

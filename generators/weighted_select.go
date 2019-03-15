package generators

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//WeightedItem is used for a generic weighted select where only a name is needed.
type WeightedItem struct {
	Name   string
	Weight int
}

type WeightedItemCollection struct {
	WeightedItem
	Items []WeightedItem
}

type WeightedBuilding struct {
	WeightedItem
	Category          string
	MaximumPercentage int
	ChildBuilding     *WeightedBuilding
	ChildBuildingName string //To be used for re-linking buildings after loading from a json file.
	ChildChance       int
	MaxQuantity       int
	MinCityWeight     int //Some structures make no sense in certain sized locations. Castle at a farm?
	MinNumEmployees   int //Bummed there is no native tuple support in Go.
	MaxNumEmployees   int
	HasOwner          bool
}

//The following block relates to loading/saving buildings to disk in the form of JSON files. While not likely, I wanted the option to create/edit buildings via text editor. We'll also use this mechanism for saving a set of default buildings in a json file vs hard coded in the project. This will also enable different building weights for different locations.

//MarshalJSON enables the encoding of WeightedBuildings to the JSON format via json.Marshal(). We perform some funky stuff to handle the fact that WeightedBuilding features a pointer. NewJSONWeightedBuilding builds a translation layer struct named JSONWeightedBuilding that stores the name of the ChildBuilding as a means to maintain this pointer. This pointer is reset when loaded from disk via the RelinkChildBuildings() function.
func (building WeightedBuilding) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONWeightedBuilding(building))

}

//UnmarshalJSON performs the inverse of MarshalJSON, dumping the provided JSON object to a translation layer struct named JSONWeightedBuilding and then using that struct to build a WeightedBuilding.
func (building *WeightedBuilding) UnmarshalJSON(data []byte) error {
	var jsonBuilding JSONWeightedBuilding
	if err := json.Unmarshal(data, &jsonBuilding); err != nil {
		return err
	}
	*building = jsonBuilding.WeightedBuilding()
	return nil

}

func NewJSONWeightedBuilding(building WeightedBuilding) JSONWeightedBuilding {
	//ChildBuilding is optional, avoid a null pointer exception by checking first.
	if building.ChildBuilding != nil {
		return JSONWeightedBuilding{building.WeightedItem.Name, building.WeightedItem.Weight, building.Category, building.MaximumPercentage, building.ChildBuilding.Name, building.ChildChance, building.MaxQuantity, building.MinCityWeight, building.MinNumEmployees, building.MaxNumEmployees, building.HasOwner}
	}
	return JSONWeightedBuilding{building.WeightedItem.Name, building.WeightedItem.Weight, building.Category, building.MaximumPercentage, "", building.ChildChance, building.MaxQuantity, building.MinCityWeight, building.MinNumEmployees, building.MaxNumEmployees, building.HasOwner}
}

//This is how we dump a WeightedBuilding to disk. The JSON keys are standardized to the JSON format and ChildBuilding becomes ChildBuildingName.
type JSONWeightedBuilding struct {
	Name              string `json:"name"`
	Weight            int    `json:"weight"`
	Category          string `json:"category"`
	MaximumPercentage int    `json:"maximumPercentage"`
	ChildBuildingName string `json:"childBuildingName"`
	ChildChance       int    `json:"childChance"`
	MaxQuantity       int    `json:"maxQuantity"`
	MinCityWeight     int    `json:"minCityWeight"`
	MinNumEmployees   int    `json:"minNumEmployee"`
	MaxNumEmployees   int    `json:"maxNumEmployee"`
	HasOwner          bool   `json:"hasOwner"`
}

//This lets us translate the translation struct to the operational struct so the project can utilize JSON encoded defaults.
func (jsonBuilding JSONWeightedBuilding) WeightedBuilding() WeightedBuilding {
	building := WeightedBuilding{}
	building.WeightedItem.Name = jsonBuilding.Name
	building.WeightedItem.Weight = jsonBuilding.Weight
	building.Category = jsonBuilding.Category
	building.MaximumPercentage = jsonBuilding.MaximumPercentage
	building.ChildBuildingName = jsonBuilding.ChildBuildingName
	building.ChildChance = jsonBuilding.ChildChance
	building.MaxQuantity = jsonBuilding.MaxQuantity
	building.MinCityWeight = jsonBuilding.MinCityWeight
	building.MinNumEmployees = jsonBuilding.MinNumEmployees
	building.MaxNumEmployees = jsonBuilding.MaxNumEmployees
	building.HasOwner = jsonBuilding.HasOwner
	return building
}

type WeightedBuildingCollection struct {
	WeightedItem
	Buildings []WeightedBuilding
}

func (item *WeightedItem) getWeight() int {
	if item.Weight == 0 {
		return 10
	}
	return item.Weight
}

func ItemsTotalWeight(items []WeightedItem) int {
	var totalWeight int
	for _, item := range items {
		totalWeight += item.getWeight()
	}
	return totalWeight
}

func RandomWeightedSelect(items []WeightedItem) (int, error) {
	totalWeight := ItemsTotalWeight(items)
	r := rand.Intn(totalWeight)
	for index, item := range items {
		r -= item.getWeight()
		if r <= 0 {
			return index, nil
		}
	}
	return -1, errors.New("no item selected")
}

//RelinkChildBuildings resets the ChildBuilding pointer in each []WeightedBuilding within a WeightedBuilding Collection. This exists as the pointer connection is broken when saving to disk. Might want to break this down so as to have a single function that relinks a building for testing purposes.
func RelinkChildBuildings(buildingCollections []*WeightedBuildingCollection) {
	for idxCollection, collection := range buildingCollections {
		for idxBuilding, building := range collection.Buildings {
			if building.ChildBuildingName != "" {
				childBuilding, err := findChild(collection, building.ChildBuildingName)
				errorHandler(err)
				//Looks like range returns a copy of the value as it iterates, so we have to operate against the original object via indices in order to set the ChildBuilding pointer.
				buildingCollections[idxCollection].Buildings[idxBuilding].ChildBuilding = childBuilding
			}
		}
	}
}

//findChild is exclusively used by RelinkChildBuildings() to rebuild the pointer link to a building's child.
func findChild(buildingCollection *WeightedBuildingCollection, buildingName string) (*WeightedBuilding, error) {
	for _, building := range buildingCollection.Buildings {
		if building.Name == buildingName {
			return &building, nil
		}
	}
	return &WeightedBuilding{}, fmt.Errorf("Unable to match building: %s", buildingName)
}

//WeightedBuildingCollectionRandomWeightedSelect is a wrapper that processes an array of buildingCollections to feed it to our generic TotalWeight and RandomWeightedSelection functions. This adds some silly wrappers in an attempt to minimize code duplication.
func WeightedBuildingCollectionRandomWeightedSelect(buildingCollections []*WeightedBuildingCollection) (int, error) {
	//My understand is initializing the slice to the known length is more efficient/faster than appending to a zero length slice.
	items := make([]WeightedItem, len(buildingCollections))
	//Extract the WeightedItems from the buildingCollections so we can do some generic selection work.
	for idx, collection := range buildingCollections {
		items[idx] = collection.WeightedItem
	}
	return RandomWeightedSelect(items)
}

func BuildingsRandomWeightedSelect(buildings []WeightedBuilding) (int, error) {
	items := make([]WeightedItem, len(buildings))
	for idx, collection := range buildings {
		items[idx] = collection.WeightedItem
	}
	return RandomWeightedSelect(items)
}

func LocationRandomWeightedSelect(locations map[string]Location) (string, error) {
	items := make([]WeightedItem, len(locations))
	names := make([]string, len(locations))
	var idx int
	for _, value := range locations {
		items[idx] = value.WeightedItem
		names[idx] = value.WeightedItem.Name
		idx++
	}
	selectedIdx, error := RandomWeightedSelect(items)
	// Using lowercase across the project for map keys.
	return strings.ToLower(names[selectedIdx]), error
}

package generators

import (
	"errors"
	"math/rand"
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
	ChildChance       int
	MaxQuantity       int
	MinCityWeight     int //Some structures make no sense in certain sized locations. Castle at a farm?
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

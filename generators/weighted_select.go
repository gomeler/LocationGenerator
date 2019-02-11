package generators

import (
	"errors"
	"math/rand"
	"time"
)

//WeightedArray in theory will form the basis for several weighted datatypes. I'm not sure if this was necessary though, might have unnecessarily complicated things.
type WeightedArray interface {
	TotalWeight() int                     //typically just used by RandomWeightedSelect
	RandomWeightedSelected() (int, error) //returns index of selected item.
}

//WeightedItem is used for a generic weighted select where only a name is needed.
type WeightedItem struct {
	Name   string
	Weight int
}

//SimpleWeightedItems is a basic collection for WeightedItems.
type SimpleWeightedItems struct {
	items []WeightedItem
}

//TotalWeight right now is solely used by RandomWeightedSelect, but I might find another use?
func (wr *SimpleWeightedItems) TotalWeight() int {
	var totalWeight int
	for _, item := range wr.items {
		totalWeight += item.Weight
	}
	return totalWeight
}

//RandomWeightedSelect is the basis for weighted random selection in this entire thing. It's kind of a big deal.
func (wr *SimpleWeightedItems) RandomWeightedSelect() (int, error) {
	totalWeight := wr.TotalWeight()
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for index, item := range wr.items {
		r -= item.Weight
		if r <= 0 {
			return index, nil
		}
	}
	return -1, errors.New("no item selected")
}

//Building related stuff. I think now that I'm just returning indices, I should have a single generic selector and type specific processors.

type WeightedBuilding struct {
	category          string
	Name              string
	MaximumPercentage int
	ChildBuilding     *WeightedBuilding
	ChildChance       int
	MaxQuantity       int
	Weight            int
	MinCityWeight     int //Some structures make no sense in certain sized locations. Castle at a farm?
}

//At some point rename this, too similar to WeightedBuilding. Maybe WeightedBuildingArray?
type WeightedBuildings struct {
	Buildings []WeightedBuilding
	Weight    int
}

func (wb *WeightedBuildings) TotalWeight() int {
	var totalWeight int
	for _, item := range wb.Buildings {
		//Some building types have no weight, we'll default to 10 for all.
		if item.Weight == 0 {
			totalWeight += 10
		} else {
		}
		totalWeight += item.Weight
	}
	return totalWeight
}

//RandomWeightedSelect is the basis for weighted random selection in this entire thing. It's kind of a big deal.
func (wb *WeightedBuildings) RandomWeightedSelect() (int, error) {
	totalWeight := wb.TotalWeight()
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for index, item := range wb.Buildings {
		//Some building types have no weight, we'll default to 10 for all.
		if item.Weight == 0 {
			r -= 10
		} else {
			r -= item.Weight
		}
		if r <= 0 {
			return index, nil
		}
	}
	return -1, errors.New("no item selected")
}

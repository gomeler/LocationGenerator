package generators

import (
	"errors"
	"fmt"
	"math"
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

type WeightedBuildings struct {
	Buildings []WeightedBuilding
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

func (wb WeightedBuilding) EvolveBuilding(maxBuildings int, buildingMap map[string]int) WeightedBuilding {
	//We receive a pointer to a WeightedBuilding with a childBuilding value. We then perform a roll to see if the childBuilding replaces the existing building. We return the child on success, otherwise we return the original value.
	var randomSelect int = rand.Intn(100)
	//var numExistingBuildings int //To be used
	if wb.ChildChance > randomSelect {
		//Success, evolution approved.
		if wb.ChildBuilding.ChildBuilding == nil {
			//No grandchild, verify there is room for the child
			if wb.ChildBuilding.MaxQuantity != 0 {
				if buildingMap[wb.ChildBuilding.Name]+1 > wb.ChildBuilding.MaxQuantity {
					//we'll have too many buildings, evolution denied.
					return wb
				}
			}
			if float64(buildingMap[wb.ChildBuilding.Name]+1) > math.Round(float64(maxBuildings*(wb.ChildBuilding.MaximumPercentage/100))) {
				//There is insufficient space.
				return wb
			}
			//Passed all the checks, child can spawn.
			fmt.Printf("Evolved %s to %s. ", wb.Name, wb.ChildBuilding.Name)
			return *wb.ChildBuilding
		} else {
			//grandchild exists, need to evolve it.
			fmt.Printf("Evolved %s to %s. ", wb.Name, wb.ChildBuilding.Name)
			return wb.ChildBuilding.EvolveBuilding(maxBuildings, buildingMap)
		}
	}
	//Failed to evolve.
	return wb
}

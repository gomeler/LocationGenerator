package generators

import (
	"errors"
	"math/rand"
	"time"
)

type WeightedItem struct {
	Name   string
	Weight int
}

func TotalWeight(items []WeightedItem) int {
	var totalWeight int
	for _, item := range items {
		totalWeight += item.Weight
	}
	return totalWeight
}

func RandomWeightedSelect(items []WeightedItem, totalWeight int) (WeightedItem, error) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for _, item := range items {
		r -= item.Weight
		if r <= 0 {
			return item, nil
		}
	}
	return WeightedItem{}, errors.New("no item selected")
}

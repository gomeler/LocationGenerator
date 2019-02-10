package data

import (
	"fmt"

	//"reflect"
	"strings"
)

type Adjectives struct {
	AdjectiveArray []string
}

func AdjMain() {
	adjData, loadError := Load(fileAdjectives)
	errorHandler(loadError)
	adjectives := new(Adjectives)
	adjectives.AdjectiveArray = strings.Split(adjData[0], ",")

	err := WriteGob(fileAdjStore, adjectives)
	errorHandler(err)

	loadedAdjectives := new(Adjectives)
	err = ReadGob(fileAdjStore, loadedAdjectives)
	errorHandler(err)
	fmt.Println(loadedAdjectives.AdjectiveArray[0])
	fmt.Println(len(loadedAdjectives.AdjectiveArray))

}

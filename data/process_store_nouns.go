package data

import (
	"fmt"

	//"reflect"
	"strings"
)

func NounMain() {
	adjData, loadError := Load(fileNouns)
	errorHandler(loadError)
	adjectives := new(Adjectives)
	adjectives.AdjectiveArray = strings.Split(adjData[0], ",")

	err := WriteGob(fileNounsStore, adjectives)
	errorHandler(err)

	loadedAdjectives := new(Adjectives)
	err = ReadGob(fileNounsStore, loadedAdjectives)
	errorHandler(err)
	fmt.Println(loadedAdjectives.AdjectiveArray[0])
	fmt.Println(len(loadedAdjectives.AdjectiveArray))

}

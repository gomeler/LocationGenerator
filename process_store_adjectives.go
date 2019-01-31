package main


import (
	"./data"
	"fmt"
	//"reflect"
	"strings"
)


type Adjectives struct {
	AdjectiveArray []string
}


const fileAdjectives string = "data/adjectives.csv"
const fileAdjStore string= "data/adjectives.gob"

func main() {
	adjData, loadError := data.Load(fileAdjectives)
	errorHandler(loadError)
	adjectives := new(Adjectives)
	adjectives.AdjectiveArray = strings.Split(adjData[0], ",")


	err := data.WriteGob(fileAdjStore, adjectives)
        errorHandler(err)


        loadedAdjectives := new (Adjectives)
        err = data.ReadGob(fileAdjStore, loadedAdjectives)
	errorHandler(err)
	fmt.Println(loadedAdjectives.AdjectiveArray[0])
	fmt.Println(len(loadedAdjectives.AdjectiveArray))



}



func errorHandler(err error) {
        if err != nil {
                panic(err)
        }
}


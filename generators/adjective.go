package generators

import (
	"math/rand"
	"time"

	"github.com/gomeler/LocationGenerator/data"
)

type Adjectives struct {
	AdjectiveArray []string
}

const fileAdjStore string = "./data/adjectives.gob"

func RandomAdjective() (string, error) {
	var adjective string
	adjectives := new(Adjectives)
	err := data.ReadGob(fileAdjStore, adjectives)
	if err == nil {
		rand.Seed(time.Now().UnixNano())
		adjective = pullAdjective(adjectives)
	}
	return adjective, err
}

func pullAdjective(adjStruct *Adjectives) string {

	var adjective string
	var idx int
	idx = rand.Intn(len(adjStruct.AdjectiveArray) - 1)
	adjective = adjStruct.AdjectiveArray[idx]
	return adjective

}

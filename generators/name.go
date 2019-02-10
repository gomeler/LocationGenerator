package generators

import (
	"math/rand"
	"time"

	"github.com/gomeler/LocationGenerator/data"
)

type Names struct {
	MaleNameArray   []string
	FemaleNameArray []string
}

const fileNames string = "./data/names.gob"

func RandomName(gender string) (string, error) {
	var name string
	names := new(Names)
	err := data.ReadGob(fileNames, names)
	if err == nil {
		//Turns out rand in go is always seeded by the value 1.
		rand.Seed(time.Now().UnixNano())
		name = pullName(gender, names)
	}
	return name, err

}

func pullName(gender string, nameStruct *Names) string {
	var name string
	var idx int
	switch gender {
	case "Male":
		idx = rand.Intn(len(nameStruct.MaleNameArray) - 1)
		name = nameStruct.MaleNameArray[idx]
	default:
		idx = rand.Intn(len(nameStruct.FemaleNameArray) - 1)
		name = nameStruct.FemaleNameArray[idx]
	}
	return name
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

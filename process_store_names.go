package main


import (
	"./data"
	"strings"
)

type Names struct {
	MaleNameArray []string
	FemaleNameArray []string
}

const fileMaleNames string = "data/male2017top.cvs"
const fileFemaleNames string = "data/female2017top.cvs"
const fileNames string = "data/names.gob"

func main() {
	maleData, maleErr := data.Load(fileMaleNames)
	femaleData, femaleErr := data.Load(fileFemaleNames)
	errorHandler(maleErr)
	errorHandler(femaleErr)
	i := 0
	maleNames := make([]string, 1)
	name := ""
	for _, value := range maleData {
		i++
		name = strings.Split(value, ",")[0]
		maleNames = append(maleNames, name)
	}

	femaleNames := make([]string, 1)
	for _, value := range femaleData {
		i++
		name = strings.Split(value, ",")[0]
		femaleNames = append(femaleNames, name)
	}

	names := Names{
		MaleNameArray: maleNames,
		FemaleNameArray: femaleNames,
	}

	err := data.WriteGob(fileNames, names)
	errorHandler(err)

	//loadedNames := new (Names)
	//err = data.ReadGob("test.gob", loadedNames)

}


func errorHandler(err error) {
        if err != nil {
                panic(err)
        }
}


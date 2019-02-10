package data

import (
	"strings"
)

type Names struct {
	MaleNameArray   []string
	FemaleNameArray []string
}

func NameMain() {
	maleData, maleErr := Load(fileMaleNames)
	femaleData, femaleErr := Load(fileFemaleNames)
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
		MaleNameArray:   maleNames,
		FemaleNameArray: femaleNames,
	}

	err := WriteGob(fileNames, names)
	errorHandler(err)

	//loadedNames := new (Names)
	//err = data.ReadGob("test.gob", loadedNames)

}

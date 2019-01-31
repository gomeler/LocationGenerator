package main

import (
	"./generators"
	"fmt"
	"strings"
)


func main() {
	race, err := generators.RandomRace()
	errorHandler(err)

	gender, err := generators.RandomGender()
	errorHandler(err)

	name, err := generators.RandomName(gender)
	errorHandler(err)
	name = strings.Replace(name, `"`, "", -1)
	name = string(name[0]) + strings.ToLower(name[1:])

	fmt.Printf("%s %s %s\n", gender, race, name)

}


func errorHandler(err error) {
        if err != nil {
                panic(err)
        }
}


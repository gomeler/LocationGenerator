package main

import (
	"fmt"
	"strings"

	"github.com/gomeler/LocationGenerator/generators"
	"github.com/gomeler/LocationGenerator/logging"
)

var log = logging.New()

func main() {
	race, err := generators.RandomRace()
	errorHandler(err)

	gender, err := generators.RandomGender()
	errorHandler(err)

	name, err := generators.RandomName(gender)
	errorHandler(err)
	name = strings.Replace(name, `"`, "", -1)
	name = string(name[0]) + strings.ToLower(name[1:])

	//holy wow, we should use log.WithFields, should increase readability a bit.
	log.Info(fmt.Sprintf("%s %s %s", gender, race, name))

}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

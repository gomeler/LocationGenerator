package generators

import (
	"fmt"
	"strings"
)

//var log = logging.New()

func CharacterEntry() {
	race, err := RandomRace()
	errorHandler(err)

	gender, err := RandomGender()
	errorHandler(err)

	name, err := RandomName(gender)
	errorHandler(err)
	name = strings.Replace(name, `"`, "", -1)
	name = string(name[0]) + strings.ToLower(name[1:])

	//holy wow, we should use log.WithFields, should increase readability a bit.
	log.Info(fmt.Sprintf("%s %s %s", gender, race, name))

}

package generators

import (
	"fmt"
	"strings"
)

//var log = logging.New()

func CharacterEntry(characterRaceFlag string, characterGenderFlag string) {
	//Right now races and genders exist primarily as a simple string. I could in theory just hand back the given flags, but there is a non-zero chance something else will happen with these facets of the generator, so I'll stick with going through the entire stack.
	race, err := RandomRace(characterRaceFlag)
	errorHandler(err)

	gender, err := RandomGender(characterGenderFlag)
	errorHandler(err)

	name, err := RandomName(gender)
	errorHandler(err)
	name = strings.Replace(name, `"`, "", -1)
	name = string(name[0]) + strings.ToLower(name[1:])

	//holy wow, we should use log.WithFields, should increase readability a bit.
	log.Info(fmt.Sprintf("%s %s %s", gender, race, name))

}

package generators

import (
	"fmt"
	"strings"
)

//NPC will form the basis of the interactive characters in a location.
type NPC struct {
	Name       string
	Gender     string
	Race       string
	Age        int
	Occupation string
}

func CharacterEntry(characterRaceFlag string, characterGenderFlag string) {
	//Moved the character gen stack to CharacterAttachment, as it needs the NPC object for attaching characters to buildings. Keeping CharacterEntry as a simple log entry to be used by the CLI for now. At some point characterCmd will be capable of processing NPC objects.
	npc := CharacterAttachment(characterRaceFlag, characterGenderFlag)
	//holy wow, we should use log.WithFields, should increase readability a bit.
	log.Info(fmt.Sprintf("%s %s %d %s", npc.Gender, npc.Race, npc.Age, npc.Name))
}

//CharacterAttachment is used to attach NPCs to buildings.
func CharacterAttachment(characterRaceFlag string, characterGenderFlag string) NPC {
	race, err := RandomRace(characterRaceFlag)
	errorHandler(err)

	gender, err := RandomGender(characterGenderFlag)
	errorHandler(err)

	name, err := RandomName(gender)
	errorHandler(err)
	name = strings.Replace(name, `"`, "", -1)
	name = string(name[0]) + strings.ToLower(name[1:])

	age := SemiNormalDistributionAgeGenerator(race)

	return NPC{Name: name, Gender: gender, Race: race.Name, Age: age}
}

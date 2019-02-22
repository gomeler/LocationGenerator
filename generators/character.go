package generators

import (
	"encoding/json"
	"fmt"
	"strings"
)

//NPC will form the basis of the interactive characters in a location.
type NPC struct {
	Name       string `json: "name"`
	Gender     string `json: "gender"`
	Race       string `json: "race"`
	Age        int    `json: "age"`
	Occupation string `json: "occupation"`
}
type NPCAlias NPC
type JSONNPC struct {
	NPCAlias
}

func NewJSONNPC(npc NPC) JSONNPC {
	return JSONNPC{NPCAlias(npc)}
}

func (jsonnpc JSONNPC) NPC() NPC {
	npc := NPC(jsonnpc.NPCAlias)
	return npc
}

//MarshalJSON exists because public/private variables in Go are stupid and encoding/json cannot see private fields in the NPC struct. Rather than just make everything public, I'm going to use this as an exercise to learn about marshal/unmarshal JSON.
func (npc NPC) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONNPC(npc))
}

func (npc *NPC) UnmarshalJSON(data []byte) error {
	var jsonNPC JSONNPC
	if err := json.Unmarshal(data, &jsonNPC); err != nil {
		return err
	}
	*npc = jsonNPC.NPC()
	return nil
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

package generators

//TODO: filename needs to be more abstract.

// locations consist of an array of min/max building quantities, to be used by a location generator.

type Location struct {
	Name              string
	MinBuildings      int
	MaxBuildings      int
	Weight            int
	GrowthFactor      float64 //To be used with events to age location.
	NPCRatio          float64 //Ratio designed to reduce the number of pre-generated NPCs. To be adjusted in the future.
	PeoplePerBuilding int     //Settling on 8 people/building unless defined in the building type.
}

var Farm = Location{"Farm", 3, 12, 1, 1.0, 0.14, 8}
var Hamlet = Location{"Hamlet", 8, 35, 2, 1.0, 0.07, 8}
var Town = Location{"Hamlet", 50, 100, 3, 1.0, 0.07, 8}
var City = Location{"City", 100, 1000, 4, 1.0, 0.04, 8}

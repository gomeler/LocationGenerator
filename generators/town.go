package generators

//TODO: filename needs to be more abstract.

// locations consist of an array of min/max building quantities, to be used by a location generator.

type Location struct {
	Name         string
	MinBuildings int
	MaxBuildings int
	Weight       int
}

var Farm = Location{"Farm", 3, 12, 1}
var Hamlet = Location{"Hamlet", 8, 35, 2}
var Town = Location{"Hamlet", 50, 100, 3}
var City = Location{"City", 100, 1000, 4}

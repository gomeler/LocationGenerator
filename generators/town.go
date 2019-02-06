package generators

//TODO: filename needs to be more abstract.

// locations consist of an array of min/max building quantities, to be used by a location generator.

type Location struct {
	Name         string
	MinBuildings int
	MaxBuildings int
	Weight       int
}

var Farm = Location{"Farm", 1, 9, 1}
var Hamlet = Location{"Hamlet", 5, 25, 2}
var Town = Location{"Hamlet", 50, 100, 3}
var City = Location{"City", 100, 1000, 4}

//var Farm = [2]int{1, 9}
//var Hamlet = [2]int{5, 25}
//var Town = [2]int{50, 100}
//var City = [2]int{100, 1000}

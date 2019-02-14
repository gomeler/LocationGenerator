package generators

const (
	House         string = "house"
	Government    string = "government"
	Medical       string = "medical"
	LightIndustry string = "lightIndustry"
	HeavyIndustry string = "heavyIndustry"
	Hospitality   string = "hospitality"
	Travel        string = "travel"
	Religious     string = "religious"
	Entertainment string = "entertainment"
)

//TODO: might also want to add a feature where if a certain building spawns, it increases the chances of related buildings spawning.

func AssembleBuildings() []*WeightedBuildingCollection {
	buildings := make([]*WeightedBuildingCollection, 9)
	buildings[0] = buildHouses()
	buildings[1] = buildMedical()
	buildings[2] = buildHospitality()
	buildings[3] = buildTravel()
	buildings[4] = buildEntertainment()
	buildings[5] = buildReligious()
	buildings[6] = buildGovernment()
	buildings[7] = buildLightIndustry()
	buildings[8] = buildHeavyIndustry()
	return buildings
}

// Assembling the housing Building array to be used for building towns.
// Houses are tiered hovel -> shack -> cottage -> townhouse.
func buildHouses() *WeightedBuildingCollection {
	var townhouse = WeightedBuilding{WeightedItem: WeightedItem{Name: "townhouse"}, Category: House, MaximumPercentage: 70, MinCityWeight: 2}
	var cottage = WeightedBuilding{WeightedItem: WeightedItem{Name: "cottage"}, Category: House, MaximumPercentage: 100, ChildBuilding: &townhouse, ChildChance: 5}
	var shack = WeightedBuilding{WeightedItem: WeightedItem{Name: "shack"}, Category: House, MaximumPercentage: 100, ChildBuilding: &cottage, ChildChance: 5}
	var hovel = WeightedBuilding{WeightedItem: WeightedItem{Name: "hovel"}, Category: House, MaximumPercentage: 100, ChildBuilding: &shack, ChildChance: 50}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: House, Weight: 140}, Buildings: []WeightedBuilding{townhouse, cottage, shack, hovel}}
}

// Medical buildings are tiered hedgeWitch -> apothecary -> hospital.
func buildMedical() *WeightedBuildingCollection {
	var hospital = WeightedBuilding{WeightedItem: WeightedItem{Name: "hospital"}, Category: Medical, MaximumPercentage: 5, MaxQuantity: 1, MinCityWeight: 3}
	var apothecary = WeightedBuilding{WeightedItem: WeightedItem{Name: "apothecary"}, Category: Medical, MaximumPercentage: 15, ChildBuilding: &hospital, ChildChance: 5, MinCityWeight: 2}
	var hedgeWitch = WeightedBuilding{WeightedItem: WeightedItem{Name: "hedgeWitch"}, Category: Medical, MaximumPercentage: 5, ChildBuilding: &apothecary, ChildChance: 20}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Medical}, Buildings: []WeightedBuilding{hospital, apothecary, hedgeWitch}}
}

// Hospitality buildings are not tiered.
func buildHospitality() *WeightedBuildingCollection {
	var tavern = WeightedBuilding{WeightedItem: WeightedItem{Name: "tavern", Weight: 10}, Category: Hospitality, MaximumPercentage: 10, MinCityWeight: 2}
	var inn = WeightedBuilding{WeightedItem: WeightedItem{Name: "inn", Weight: 10}, Category: Hospitality, MaximumPercentage: 10, MinCityWeight: 2}
	var hostel = WeightedBuilding{WeightedItem: WeightedItem{Name: "hostel", Weight: 10}, Category: Hospitality, MaximumPercentage: 10, MinCityWeight: 2}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Hospitality, Weight: 20}, Buildings: []WeightedBuilding{tavern, inn, hostel}}
}

// Travel buildings feature two branches, coast -> lighthouse, and river -> bridge.
func buildTravel() *WeightedBuildingCollection {
	var lighthouse = WeightedBuilding{WeightedItem: WeightedItem{Name: "lighthouse"}, Category: Travel, MaximumPercentage: 10, MaxQuantity: 1}
	var coast = WeightedBuilding{WeightedItem: WeightedItem{Name: "coast"}, Category: Travel, MaximumPercentage: 10, ChildBuilding: &lighthouse, ChildChance: 25, MaxQuantity: 1}
	var bridge = WeightedBuilding{WeightedItem: WeightedItem{Name: "bridge"}, Category: Travel, MaximumPercentage: 10}
	var river = WeightedBuilding{WeightedItem: WeightedItem{Name: "river"}, Category: Travel, MaximumPercentage: 10, ChildBuilding: &bridge, ChildChance: 25}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Travel, Weight: 20}, Buildings: []WeightedBuilding{lighthouse, coast, bridge, river}}
}

// Entertainment buildings feature an arena -> stadium -> ampitheatre tier, and an untiered operaHouse.
func buildEntertainment() *WeightedBuildingCollection {
	var ampitheatre = WeightedBuilding{WeightedItem: WeightedItem{Name: "ampitheatre"}, Category: Entertainment, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 4}
	var stadium = WeightedBuilding{WeightedItem: WeightedItem{Name: "stadium"}, Category: Entertainment, MaximumPercentage: 10, ChildBuilding: &ampitheatre, ChildChance: 25, MinCityWeight: 3}
	var arena = WeightedBuilding{WeightedItem: WeightedItem{Name: "arena"}, Category: Entertainment, MaximumPercentage: 10, ChildBuilding: &stadium, ChildChance: 25, MinCityWeight: 2}
	var operaHouse = WeightedBuilding{WeightedItem: WeightedItem{Name: "opera house"}, Category: Entertainment, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Entertainment, Weight: 10}, Buildings: []WeightedBuilding{ampitheatre, stadium, arena, operaHouse}}
}

// Religious buildings feature a shrine -> temple tier, and a (nunnery, abbey) -> church -> cathedral tier.
func buildReligious() *WeightedBuildingCollection {
	var cathedral = WeightedBuilding{WeightedItem: WeightedItem{Name: "cathedral"}, Category: Religious, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	var church = WeightedBuilding{WeightedItem: WeightedItem{Name: "church"}, Category: Religious, MaximumPercentage: 10, ChildBuilding: &cathedral, ChildChance: 10, MinCityWeight: 2}
	var abbey = WeightedBuilding{WeightedItem: WeightedItem{Name: "abbey"}, Category: Religious, MaximumPercentage: 10, ChildBuilding: &church, ChildChance: 10}
	var nunnery = WeightedBuilding{WeightedItem: WeightedItem{Name: "nunnery"}, Category: Religious, MaximumPercentage: 10, ChildBuilding: &church, ChildChance: 10}
	var temple = WeightedBuilding{WeightedItem: WeightedItem{Name: "temple"}, Category: Religious, MaximumPercentage: 10, MinCityWeight: 2}
	var shrine = WeightedBuilding{WeightedItem: WeightedItem{Name: "shrine"}, Category: Religious, MaximumPercentage: 10, ChildBuilding: &temple, ChildChance: 10}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Religious, Weight: 15}, Buildings: []WeightedBuilding{cathedral, church, abbey, nunnery, temple, shrine}}
}

// Government buildings have a few tiers. library, walls, and palace have no children. (mayorHouse -> assembly building, moot building) -> townhall. guardOutpost -> fortress -> castle. I figure a town can have walls without necessarily having a guard outpost. Jail -> prison seems logical.
func buildGovernment() *WeightedBuildingCollection {
	var library = WeightedBuilding{WeightedItem: WeightedItem{Name: "library"}, Category: Government, MaximumPercentage: 10, MinCityWeight: 3}
	var townhall = WeightedBuilding{WeightedItem: WeightedItem{Name: "townhall"}, Category: Government, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	var mootBuilding = WeightedBuilding{WeightedItem: WeightedItem{Name: "moot building"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &townhall, ChildChance: 20, MaxQuantity: 1, MinCityWeight: 3}
	var assemblyBuilding = WeightedBuilding{WeightedItem: WeightedItem{Name: "assembly building"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &townhall, ChildChance: 20, MaxQuantity: 1, MinCityWeight: 2}
	var mayorHouse = WeightedBuilding{WeightedItem: WeightedItem{Name: "mayor's house"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &assemblyBuilding, ChildChance: 10, MaxQuantity: 1, MinCityWeight: 1}
	var castle = WeightedBuilding{WeightedItem: WeightedItem{Name: "castle"}, Category: Government, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	var fortress = WeightedBuilding{WeightedItem: WeightedItem{Name: "fortress"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &castle, ChildChance: 20, MaxQuantity: 1, MinCityWeight: 3}
	var guardOutpost = WeightedBuilding{WeightedItem: WeightedItem{Name: "guard outpost"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &fortress, ChildChance: 20, MinCityWeight: 2}
	var prison = WeightedBuilding{WeightedItem: WeightedItem{Name: "prison"}, Category: Government, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	var jail = WeightedBuilding{WeightedItem: WeightedItem{Name: "jail"}, Category: Government, MaximumPercentage: 10, ChildBuilding: &prison, ChildChance: 20, MinCityWeight: 2}
	var palace = WeightedBuilding{WeightedItem: WeightedItem{Name: "palace"}, Category: Government, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	var walls = WeightedBuilding{WeightedItem: WeightedItem{Name: "walls"}, Category: Government, MaximumPercentage: 10, MaxQuantity: 1, MinCityWeight: 3}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: Government, Weight: 15}, Buildings: []WeightedBuilding{library, townhall, mootBuilding, assemblyBuilding, mayorHouse, fortress, castle, guardOutpost, jail, prison, palace, walls}}
}

//The industry buildings are a bit of a mess. I think I'd like some relationship between the light and heavy industry buildings, but for now I'll keep them as two untiered arrays.
func buildLightIndustry() *WeightedBuildingCollection {
	var pier = WeightedBuilding{WeightedItem: WeightedItem{Name: "pier"}, Category: LightIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var blacksmith = WeightedBuilding{WeightedItem: WeightedItem{Name: "blacksmith"}, Category: LightIndustry, MaximumPercentage: 10}
	var leatherworker = WeightedBuilding{WeightedItem: WeightedItem{Name: "leather worker"}, Category: LightIndustry, MaximumPercentage: 10}
	var butchershop = WeightedBuilding{WeightedItem: WeightedItem{Name: "butchershop"}, Category: LightIndustry, MaximumPercentage: 10}
	var bakery = WeightedBuilding{WeightedItem: WeightedItem{Name: "bakery"}, Category: LightIndustry, MaximumPercentage: 10}
	var weaver = WeightedBuilding{WeightedItem: WeightedItem{Name: "weaver"}, Category: LightIndustry, MaximumPercentage: 10}
	var machinist = WeightedBuilding{WeightedItem: WeightedItem{Name: "machinist"}, Category: LightIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var mason = WeightedBuilding{WeightedItem: WeightedItem{Name: "mason"}, Category: LightIndustry, MaximumPercentage: 10}
	var woodworker = WeightedBuilding{WeightedItem: WeightedItem{Name: "woodworker"}, Category: LightIndustry, MaximumPercentage: 10}
	var generalstore = WeightedBuilding{WeightedItem: WeightedItem{Name: "general store"}, Category: LightIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: LightIndustry, Weight: 25}, Buildings: []WeightedBuilding{pier, blacksmith, leatherworker, butchershop, bakery, weaver, machinist, mason, woodworker, generalstore}}
}

//Related to the light building comment, I'd like to create some links between light and heavy industry buildings. Some of them naturally make sense as pairs, but I also wanted to separate them due to some of the heavy industry buildings not making sense in a town that calls for only light industry. Might just have to combine them with a value that defines the minimum town size?
func buildHeavyIndustry() *WeightedBuildingCollection {
	var dock = WeightedBuilding{WeightedItem: WeightedItem{Name: "dock"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 3}
	var harbor = WeightedBuilding{WeightedItem: WeightedItem{Name: "harbor"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 3}
	var shipbuilder = WeightedBuilding{WeightedItem: WeightedItem{Name: "shipbuilder"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 3}
	var brickmaker = WeightedBuilding{WeightedItem: WeightedItem{Name: "brickmaker"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var forge = WeightedBuilding{WeightedItem: WeightedItem{Name: "forge"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var furnace = WeightedBuilding{WeightedItem: WeightedItem{Name: "furnace"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var warehouse = WeightedBuilding{WeightedItem: WeightedItem{Name: "warehouse"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 3}
	var surfaceMine = WeightedBuilding{WeightedItem: WeightedItem{Name: "surface mine"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	var subsurfaceMine = WeightedBuilding{WeightedItem: WeightedItem{Name: "sub-surface mine"}, Category: HeavyIndustry, MaximumPercentage: 10, MinCityWeight: 2}
	return &WeightedBuildingCollection{WeightedItem: WeightedItem{Name: HeavyIndustry, Weight: 10}, Buildings: []WeightedBuilding{dock, harbor, shipbuilder, brickmaker, forge, furnace, warehouse, surfaceMine, subsurfaceMine}}
}

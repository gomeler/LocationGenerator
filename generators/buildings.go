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

func AssembleBuildings() map[string]*WeightedBuildings {
	var buildings = make(map[string]*WeightedBuildings)
	buildings[House] = buildHouses()
	buildings[Medical] = buildMedical()
	buildings[Hospitality] = buildHospitality()
	buildings[Travel] = buildTravel()
	buildings[Entertainment] = buildEntertainment()
	buildings[Religious] = buildReligious()
	buildings[Government] = buildGovernment()
	buildings[LightIndustry] = buildLightIndustry()
	buildings[HeavyIndustry] = buildHeavyIndustry()
	return buildings
}

// Assembling the housing Building array to be used for building towns.
// Houses are tiered hovel -> shack -> cottage -> townhouse.
func buildHouses() *WeightedBuildings {
	var townhouse = WeightedBuilding{category: House, MaximumPercentage: 70, Name: "townhouse"}
	var cottage = WeightedBuilding{category: House, Name: "cottage", MaximumPercentage: 100, ChildBuilding: &townhouse, ChildChance: 5}
	var shack = WeightedBuilding{category: House, Name: "shack", MaximumPercentage: 100, ChildBuilding: &cottage, ChildChance: 5}
	var hovel = WeightedBuilding{category: House, Name: "hovel", MaximumPercentage: 100, ChildBuilding: &shack, ChildChance: 50}
	return &WeightedBuildings{Buildings: []WeightedBuilding{townhouse, cottage, shack, hovel}}
}

// Medical buildings are tiered hedgeWitch -> apothecary -> hospital.
func buildMedical() *WeightedBuildings {
	var hospital = WeightedBuilding{category: Medical, Name: "hospital", MaximumPercentage: 5, MaxQuantity: 1}
	var apothecary = WeightedBuilding{category: Medical, Name: "apothecary", MaximumPercentage: 15, ChildBuilding: &hospital, ChildChance: 5}
	var hedgeWitch = WeightedBuilding{category: Medical, Name: "hedgeWitch", MaximumPercentage: 5, ChildBuilding: &apothecary, ChildChance: 20}
	return &WeightedBuildings{Buildings: []WeightedBuilding{hospital, apothecary, hedgeWitch}}
}

// Hospitality buildings are not tiered.
func buildHospitality() *WeightedBuildings {
	var tavern = WeightedBuilding{category: Hospitality, Name: "tavern", MaximumPercentage: 10, Weight: 10}
	var inn = WeightedBuilding{category: Hospitality, Name: "inn", MaximumPercentage: 10, Weight: 10}
	var hostel = WeightedBuilding{category: Hospitality, Name: "hostel", MaximumPercentage: 10, Weight: 10}
	return &WeightedBuildings{Buildings: []WeightedBuilding{tavern, inn, hostel}}
}

// Travel buildings feature two branches, coast -> lighthouse, and river -> bridge.
func buildTravel() *WeightedBuildings {
	var lighthouse = WeightedBuilding{category: Travel, Name: "lighthouse", MaximumPercentage: 10, MaxQuantity: 1}
	var coast = WeightedBuilding{category: Travel, Name: "coast", MaximumPercentage: 10, ChildBuilding: &lighthouse, ChildChance: 25, MaxQuantity: 1}
	var bridge = WeightedBuilding{category: Travel, Name: "bridge", MaximumPercentage: 10}
	var river = WeightedBuilding{category: Travel, Name: "river", MaximumPercentage: 10, ChildBuilding: &bridge, ChildChance: 25}
	return &WeightedBuildings{Buildings: []WeightedBuilding{lighthouse, coast, bridge, river}}
}

// Entertainment buildings feature an arena -> stadium -> ampitheatre tier, and an untiered operaHouse.
func buildEntertainment() *WeightedBuildings {
	var ampitheatre = WeightedBuilding{category: Entertainment, Name: "ampitheatre", MaximumPercentage: 10, MaxQuantity: 1}
	var stadium = WeightedBuilding{category: Entertainment, Name: "stadium", MaximumPercentage: 10, ChildBuilding: &ampitheatre, ChildChance: 25}
	var arena = WeightedBuilding{category: Entertainment, Name: "arena", MaximumPercentage: 10, ChildBuilding: &stadium, ChildChance: 25}
	var operaHouse = WeightedBuilding{category: Entertainment, Name: "opera house", MaximumPercentage: 10, MaxQuantity: 1}
	return &WeightedBuildings{Buildings: []WeightedBuilding{ampitheatre, stadium, arena, operaHouse}}
}

// Religious buildings feature a shrine -> temple tier, and a (nunnery, abbey) -> church -> cathedral tier.
func buildReligious() *WeightedBuildings {
	var cathedral = WeightedBuilding{category: Religious, Name: "cathedral", MaximumPercentage: 10, MaxQuantity: 1}
	var church = WeightedBuilding{category: Religious, Name: "church", MaximumPercentage: 10, ChildBuilding: &cathedral, ChildChance: 10}
	var abbey = WeightedBuilding{category: Religious, Name: "abbey", MaximumPercentage: 10, ChildBuilding: &church, ChildChance: 10}
	var nunnery = WeightedBuilding{category: Religious, Name: "nunnery", MaximumPercentage: 10, ChildBuilding: &church, ChildChance: 10}
	var temple = WeightedBuilding{category: Religious, Name: "temple", MaximumPercentage: 10}
	var shrine = WeightedBuilding{category: Religious, Name: "shrine", MaximumPercentage: 10, ChildBuilding: &temple, ChildChance: 10}
	return &WeightedBuildings{Buildings: []WeightedBuilding{cathedral, church, abbey, nunnery, temple, shrine}}
}

// Government buildings have a few tiers. library, walls, and palace have no children. (mayorHouse -> assembly building, moot building) -> townhall. guardOutpost -> fortress -> castle. I figure a town can have walls without necessarily having a guard outpost. Jail -> prison seems logical.
func buildGovernment() *WeightedBuildings {
	var library = WeightedBuilding{category: Government, Name: "library", MaximumPercentage: 10}
	var townhall = WeightedBuilding{category: Government, Name: "townhall", MaximumPercentage: 10, MaxQuantity: 1}
	var mootBuilding = WeightedBuilding{category: Government, Name: "moot building", MaximumPercentage: 10, ChildBuilding: &townhall, ChildChance: 20, MaxQuantity: 1}
	var assemblyBuilding = WeightedBuilding{category: Government, Name: "assembly building", MaximumPercentage: 10, ChildBuilding: &townhall, ChildChance: 20, MaxQuantity: 1}
	var mayorHouse = WeightedBuilding{category: Government, Name: "mayor's house", MaximumPercentage: 10, ChildBuilding: &assemblyBuilding, ChildChance: 10, MaxQuantity: 1}
	var castle = WeightedBuilding{category: Government, Name: "castle", MaximumPercentage: 10, MaxQuantity: 1}
	var fortress = WeightedBuilding{category: Government, Name: "fortress", MaximumPercentage: 10, ChildBuilding: &castle, ChildChance: 20, MaxQuantity: 1}
	var guardOutpost = WeightedBuilding{category: Government, Name: "guard outpost", MaximumPercentage: 10, ChildBuilding: &fortress, ChildChance: 20}
	var prison = WeightedBuilding{category: Government, Name: "prison", MaximumPercentage: 10, MaxQuantity: 1}
	var jail = WeightedBuilding{category: Government, Name: "jail", MaximumPercentage: 10, ChildBuilding: &prison, ChildChance: 20}
	var palace = WeightedBuilding{category: Government, Name: "palace", MaximumPercentage: 10, MaxQuantity: 1}
	var walls = WeightedBuilding{category: Government, Name: "walls", MaximumPercentage: 10, MaxQuantity: 1}
	return &WeightedBuildings{Buildings: []WeightedBuilding{library, townhall, mootBuilding, assemblyBuilding, mayorHouse, fortress, castle, guardOutpost, jail, prison, palace, walls}}
}

//The industry buildings are a bit of a mess. I think I'd like some relationship between the light and heavy industry buildings, but for now I'll keep them as two untiered arrays.
func buildLightIndustry() *WeightedBuildings {
	var pier = WeightedBuilding{category: LightIndustry, Name: "pier", MaximumPercentage: 10}
	var blacksmith = WeightedBuilding{category: LightIndustry, Name: "blacksmith", MaximumPercentage: 10}
	var leatherworker = WeightedBuilding{category: LightIndustry, Name: "leather worker", MaximumPercentage: 10}
	var butchershop = WeightedBuilding{category: LightIndustry, Name: "butchershop", MaximumPercentage: 10}
	var bakery = WeightedBuilding{category: LightIndustry, Name: "bakery", MaximumPercentage: 10}
	var weaver = WeightedBuilding{category: LightIndustry, Name: "weaver", MaximumPercentage: 10}
	var machinist = WeightedBuilding{category: LightIndustry, Name: "machinist", MaximumPercentage: 10}
	var mason = WeightedBuilding{category: LightIndustry, Name: "mason", MaximumPercentage: 10}
	var woodworker = WeightedBuilding{category: LightIndustry, Name: "woodworker", MaximumPercentage: 10}
	var generalstore = WeightedBuilding{category: LightIndustry, Name: "general store", MaximumPercentage: 10}
	return &WeightedBuildings{Buildings: []WeightedBuilding{pier, blacksmith, leatherworker, butchershop, bakery, weaver, machinist, mason, woodworker, generalstore}}
}

//Related to the light building comment, I'd like to create some links between light and heavy industry buildings. Some of them naturally make sense as pairs, but I also wanted to separate them due to some of the heavy industry buildings not making sense in a town that calls for only light industry. Might just have to combine them with a value that defines the minimum town size?
func buildHeavyIndustry() *WeightedBuildings {
	var dock = WeightedBuilding{category: HeavyIndustry, Name: "dock", MaximumPercentage: 10}
	var harbor = WeightedBuilding{category: HeavyIndustry, Name: "harbor", MaximumPercentage: 10}
	var shipbuilder = WeightedBuilding{category: HeavyIndustry, Name: "shipbuilder", MaximumPercentage: 10}
	var brickmaker = WeightedBuilding{category: HeavyIndustry, Name: "brickmaker", MaximumPercentage: 10}
	var forge = WeightedBuilding{category: HeavyIndustry, Name: "forge", MaximumPercentage: 10}
	var furnace = WeightedBuilding{category: HeavyIndustry, Name: "furnace", MaximumPercentage: 10}
	var warehouse = WeightedBuilding{category: HeavyIndustry, Name: "warehouse", MaximumPercentage: 10}
	var surfaceMine = WeightedBuilding{category: HeavyIndustry, Name: "surface mine", MaximumPercentage: 10}
	var subsurfaceMine = WeightedBuilding{category: HeavyIndustry, Name: "sub-surface mine", MaximumPercentage: 10}
	return &WeightedBuildings{Buildings: []WeightedBuilding{dock, harbor, shipbuilder, brickmaker, forge, furnace, warehouse, surfaceMine, subsurfaceMine}}
}

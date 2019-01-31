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

/* Buildings are all based on the simple struct. The important link is the childBuilding pointer. Every category of building has a base/parent, and each parent has a chance of evolving to a child Building. For each category of building, all of the buildings will be evaluated for each child/parent relationship.*/

//TODO: might also want to add a feature where if a certain building spawns, it increases the chances of related buildings spawning.
type Building struct {
	category          string
	Name              string
	maximumPercentage int
	childBuilding     *Building
	childChance       int
	maxQuantity       int
	WeightedItem
}

func AssembleBuildings() map[string][]Building {
	var buildings = make(map[string][]Building)
	//var houseBuildings = buildHouses()
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
func buildHouses() []Building {
	var townhouse = Building{category: House, maximumPercentage: 70, Name: "townhouse"}
	var cottage = Building{category: House, Name: "cottage", maximumPercentage: 100, childBuilding: &townhouse, childChance: 5}
	var shack = Building{category: House, Name: "shack", maximumPercentage: 100, childBuilding: &cottage, childChance: 5}
	var hovel = Building{category: House, Name: "hovel", maximumPercentage: 100, childBuilding: &shack, childChance: 50}
	return []Building{townhouse, cottage, shack, hovel}
}

// Medical buildings are tiered hedgeWitch -> apothecary -> hospital.
func buildMedical() []Building {
	var hospital = Building{category: Medical, Name: "hospital", maximumPercentage: 5, maxQuantity: 1}
	var apothecary = Building{category: Medical, Name: "apothecary", maximumPercentage: 15, childBuilding: &hospital, childChance: 5}
	var hedgeWitch = Building{category: Medical, Name: "hedgeWitch", maximumPercentage: 5, childBuilding: &apothecary, childChance: 20}
	return []Building{hospital, apothecary, hedgeWitch}
}

// Hospitality buildings are not tiered.
func buildHospitality() []Building {
	var tavern = Building{category: Hospitality, maximumPercentage: 10, WeightedItem: WeightedItem{Name: "tavern", Weight: 10}}
	var inn = Building{category: Hospitality, maximumPercentage: 10, WeightedItem: WeightedItem{Name: "inn", Weight: 15}}
	var hostel = Building{category: Hospitality, maximumPercentage: 10, WeightedItem: WeightedItem{Name: "hostel", Weight: 5}}
	return []Building{tavern, inn, hostel}
}

// Travel buildings feature two branches, coast -> lighthouse, and river -> bridge.
func buildTravel() []Building {
	var lighthouse = Building{category: Travel, Name: "lighthouse", maximumPercentage: 10, maxQuantity: 1}
	var coast = Building{category: Travel, Name: "coast", maximumPercentage: 10, childBuilding: &lighthouse, childChance: 25}
	var bridge = Building{category: Travel, Name: "bridge", maximumPercentage: 10}
	var river = Building{category: Travel, Name: "river", maximumPercentage: 10, childBuilding: &bridge, childChance: 25}
	return []Building{lighthouse, coast, bridge, river}
}

// Entertainment buildings feature an arena -> stadium -> ampitheatre tier, and an untiered operaHouse.
func buildEntertainment() []Building {
	var ampitheatre = Building{category: Entertainment, Name: "ampitheatre", maximumPercentage: 10, maxQuantity: 1}
	var stadium = Building{category: Entertainment, Name: "stadium", maximumPercentage: 10, childBuilding: &ampitheatre, childChance: 25}
	var arena = Building{category: Entertainment, Name: "arena", maximumPercentage: 10, childBuilding: &stadium, childChance: 25}
	var operaHouse = Building{category: Entertainment, Name: "opera house", maximumPercentage: 10, maxQuantity: 1}
	return []Building{ampitheatre, stadium, arena, operaHouse}
}

// Religious buildings feature a shrine -> temple tier, and a (nunnery, abbey) -> church -> cathedral tier.
func buildReligious() []Building {
	var cathedral = Building{category: Religious, Name: "cathedral", maximumPercentage: 10, maxQuantity: 1}
	var church = Building{category: Religious, Name: "church", maximumPercentage: 10, childBuilding: &cathedral, childChance: 10}
	var abbey = Building{category: Religious, Name: "abbey", maximumPercentage: 10, childBuilding: &church, childChance: 10}
	var nunnery = Building{category: Religious, Name: "nunnery", maximumPercentage: 10, childBuilding: &church, childChance: 10}
	var temple = Building{category: Religious, Name: "temple", maximumPercentage: 10}
	var shrine = Building{category: Religious, Name: "shrine", maximumPercentage: 10, childBuilding: &temple, childChance: 10}
	return []Building{cathedral, church, abbey, nunnery, temple, shrine}
}

// Government buildings have a few tiers. library, walls, and palace have no children. (mayorHouse -> assembly building, moot building) -> townhall. guardOutpost -> fortress -> castle. I figure a town can have walls without necessarily having a guard outpost. Jail -> prison seems logical.
func buildGovernment() []Building {
	var library = Building{category: Government, Name: "library", maximumPercentage: 10}
	var townhall = Building{category: Government, Name: "townhall", maximumPercentage: 10, maxQuantity: 1}
	var mootBuilding = Building{category: Government, Name: "moot building", maximumPercentage: 10, childBuilding: &townhall, childChance: 20, maxQuantity: 1}
	var assemblyBuilding = Building{category: Government, Name: "assembly building", maximumPercentage: 10, childBuilding: &townhall, childChance: 20, maxQuantity: 1}
	var mayorHouse = Building{category: Government, Name: "mayor's house", maximumPercentage: 10, childBuilding: &assemblyBuilding, childChance: 10, maxQuantity: 1}
	var castle = Building{category: Government, Name: "castle", maximumPercentage: 10, maxQuantity: 1}
	var fortress = Building{category: Government, Name: "fortress", maximumPercentage: 10, childBuilding: &castle, childChance: 20, maxQuantity: 1}
	var guardOutpost = Building{category: Government, Name: "guard outpost", maximumPercentage: 10, childBuilding: &fortress, childChance: 20}
	var prison = Building{category: Government, Name: "prison", maximumPercentage: 10, maxQuantity: 1}
	var jail = Building{category: Government, Name: "jail", maximumPercentage: 10, childBuilding: &prison, childChance: 20}
	var palace = Building{category: Government, Name: "palace", maximumPercentage: 10, maxQuantity: 1}
	var walls = Building{category: Government, Name: "walls", maximumPercentage: 10, maxQuantity: 1}
	return []Building{library, townhall, mootBuilding, assemblyBuilding, mayorHouse, fortress, castle, guardOutpost, jail, prison, palace, walls}
}

//The industry buildings are a bit of a mess. I think I'd like some relationship between the light and heavy industry buildings, but for now I'll keep them as two untiered arrays.
func buildLightIndustry() []Building {
	var pier = Building{category: LightIndustry, Name: "pier", maximumPercentage: 10}
	var blacksmith = Building{category: LightIndustry, Name: "blacksmith", maximumPercentage: 10}
	var leatherworker = Building{category: LightIndustry, Name: "leather worker", maximumPercentage: 10}
	var butchershop = Building{category: LightIndustry, Name: "butchershop", maximumPercentage: 10}
	var bakery = Building{category: LightIndustry, Name: "bakery", maximumPercentage: 10}
	var weaver = Building{category: LightIndustry, Name: "weaver", maximumPercentage: 10}
	var machinist = Building{category: LightIndustry, Name: "machinist", maximumPercentage: 10}
	var mason = Building{category: LightIndustry, Name: "mason", maximumPercentage: 10}
	var woodworker = Building{category: LightIndustry, Name: "woodworker", maximumPercentage: 10}
	var generalstore = Building{category: LightIndustry, Name: "general store", maximumPercentage: 10}
	return []Building{pier, blacksmith, leatherworker, butchershop, bakery, weaver, machinist, mason, woodworker, generalstore}
}

//Related to the light building comment, I'd like to create some links between light and heavy industry buildings. Some of them naturally make sense as pairs, but I also wanted to separate them due to some of the heavy industry buildings not making sense in a town that calls for only light industry. Might just have to combine them with a value that defines the minimum town size?
func buildHeavyIndustry() []Building {
	var dock = Building{category: HeavyIndustry, Name: "dock", maximumPercentage: 10}
	var harbor = Building{category: HeavyIndustry, Name: "harbor", maximumPercentage: 10}
	var shipbuilder = Building{category: HeavyIndustry, Name: "shipbuilder", maximumPercentage: 10}
	var brickmaker = Building{category: HeavyIndustry, Name: "brickmaker", maximumPercentage: 10}
	var forge = Building{category: HeavyIndustry, Name: "forge", maximumPercentage: 10}
	var furnace = Building{category: HeavyIndustry, Name: "furnace", maximumPercentage: 10}
	var warehouse = Building{category: HeavyIndustry, Name: "warehouse", maximumPercentage: 10}
	var surfaceMine = Building{category: HeavyIndustry, Name: "surface mine", maximumPercentage: 10}
	var subsurfaceMine = Building{category: HeavyIndustry, Name: "sub-surface mine", maximumPercentage: 10}
	return []Building{dock, harbor, shipbuilder, brickmaker, forge, furnace, warehouse, surfaceMine, subsurfaceMine}
}

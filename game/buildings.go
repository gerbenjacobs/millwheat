package game

const (
	BuildingFarm BuildingType = iota
	BuildingMill
	BuildingBakery
)

type Buildings map[BuildingType]Building

type BuildingType int

type TownBuilding struct {
	Type         BuildingType
	CurrentLevel int
}

type Building struct {
	Name        string
	Description string
	Image       string
	Consumes    []ItemID
	Produces    []ItemID
	Mechanics   []BuildingMechanic
	BuildCosts  map[int]BuildCost
}

// BuildingMechanic contains the mechanical and proficiency attributes of buildings
// ex. Farm: Wheat per hour
type BuildingMechanic struct {
	Name   string
	Levels map[int]int
}

type BuildCost struct {
	Stones int
	Planks int
}

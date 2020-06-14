package game

import (
	"github.com/google/uuid"
)

const (
	BuildingFarm BuildingType = iota
	BuildingMill
	BuildingBakery
)

const ( // const blocks separated to reset iota
	MechanicConsumption MechanicType = iota
	MechanicEfficiency
	MechanicOutput
)

type Buildings map[BuildingType]Building

type BuildingType int
type MechanicType int

// TownBuilding represents an instance of a building, located in a town
type TownBuilding struct {
	ID           uuid.UUID
	Type         BuildingType
	CurrentLevel int
}

// Building contains all the data for buildings in the game
type Building struct {
	Name        string
	Description string
	Image       string
	Consumes    []ItemID
	Produces    []ItemID
	IsGenerator bool
	Mechanics   []BuildingMechanic
	BuildCosts  map[int]BuildCost
}

// BuildingMechanic contains the mechanical and proficiency attributes of buildings
// ex. Farm: Wheat per hour
type BuildingMechanic struct {
	Type   MechanicType
	Name   string
	Levels map[int]int
}

type BuildCost struct {
	Stones int
	Planks int
}

func (m MechanicType) String() string {
	return []string{"Consumption", "Efficiency", "Output"}[m]
}

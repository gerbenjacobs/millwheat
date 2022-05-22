package game

import (
	"sort"
)

const (
	MechanicConsumption MechanicType = iota
	MechanicEfficiency
	MechanicOutput
)

type MechanicType int

// BuildingMechanic contains the mechanical and proficiency attributes of buildings
// ex. Farm: Wheat per hour
type BuildingMechanic struct {
	Type   MechanicType
	Name   string
	ItemID ItemID
	Levels map[int]int
}

func (m MechanicType) String() string {
	return []string{"Consumption", "Efficiency", "Output"}[m]
}

func (b Building) MechanicsList() []BuildingMechanic {
	var mList []BuildingMechanic
	for _, i := range b.Mechanics {
		mList = append(mList, i)
	}
	sort.SliceStable(b.Mechanics, func(i, j int) bool {
		if mList[i].Type == mList[j].Type {
			return mList[i].Name < mList[j].Name
		}
		return mList[i].Type < mList[j].Type
	})
	return mList
}

func (b Building) MaxConsumption(itemID ItemID, level int) int {
	for _, m := range b.Mechanics {
		if m.Type == MechanicConsumption && m.ItemID == itemID {
			return m.Levels[level]
		}
	}
	return 0
}

func (b Building) MaxProduction(itemID ItemID, level int) int {
	for _, m := range b.Mechanics {
		if m.Type == MechanicOutput && m.ItemID == itemID {
			return m.Levels[level]
		}
	}
	// if no output found, see if we can determine based on consumption
	// this only applies to buildings where the key of production is consumption=true
	// at the moment that's only the Butcher
	for production := range b.Production {
		if production.ItemID == itemID && production.IsConsumption {
			return b.MaxConsumption(itemID, level)
		}
	}
	return 0
}

func (b Building) MaxEfficiency(itemID ItemID, level int) int {
	for _, m := range b.Mechanics {
		if m.Type == MechanicEfficiency && m.ItemID == itemID {
			return m.Levels[level]
		}
	}
	return 0
}

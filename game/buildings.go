package game

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

const (
	BuildingFarm BuildingType = iota
	BuildingMill
	BuildingBakery
	BuildingPigFarm
	BuildingButcher
	BuildingWeaponSmith
)

type Buildings map[BuildingType]Building
type BuildingType int

// TownBuilding represents an instance of a building, located in a town
type TownBuilding struct {
	ID           uuid.UUID
	Type         BuildingType
	CurrentLevel int
}

// Building contains all the data for buildings in the game
type Building struct {
	Name           string
	Description    string
	Image          string
	Production     map[ItemSet][]ItemSet
	IsGenerator    bool
	LastCollection time.Time
	Mechanics      []BuildingMechanic
	BuildCosts     map[int]BuildingCost
}

// BuildingCost contains the cost in stones and planks
type BuildingCost struct {
	Stones int
	Planks int
}

// ItemSet is used to calculate what items you need to produce something
type ItemSet struct {
	ItemID        ItemID
	Quantity      int
	IsConsumption bool
}

// ProductionResult contains calculated values for producing an item
type ProductionResult struct {
	Consumption []ItemSet
	Production  []ItemSet
	Hours       int
}

func (b Building) CanDealWith(id ItemID) bool {
	for _, itemID := range append(b.ConsumesList(), b.ProducesList()...) {
		if id == itemID {
			return true
		}
	}

	return false
}

func (b Building) CreateProduct(product ItemID, quantity, level int) (*ProductionResult, error) {
	var consume []ItemSet
	var produce []ItemSet
	var isConsumable = false
	for item, subItems := range b.Production {
		if item.ItemID == product {
			if item.IsConsumption {
				isConsumable = true
				consume = append(consume, item)
			} else {
				produce = append(produce, item)
			}

			for _, subItem := range subItems {
				if subItem.IsConsumption {
					consume = append(consume, subItem)
				} else {
					produce = append(produce, subItem)
				}
			}
		}
	}

	if len(consume) == 0 || len(produce) == 0 {
		return nil, errors.New("building can not create this item")
	}

	// Calculate consumption and production via efficiency along with estimated time in hours
	for i, c := range consume {
		e := b.MaxEfficiency(c.ItemID, level)
		if e == 0 {
			consume[i].Quantity = quantity
		} else {
			// consume less; divide
			consume[i].Quantity = int(math.Ceil(float64(quantity) / float64(e)))
		}
	}
	for i, p := range produce {
		e := b.MaxEfficiency(p.ItemID, level)
		if e == 0 {
			produce[i].Quantity = quantity
		} else {
			// produce more; multiply
			produce[i].Quantity = quantity * e
		}
	}

	var div int
	if isConsumable {
		div = b.MaxConsumption(product, level)
	} else {
		div = b.MaxProduction(product, level)
	}
	hours := math.Ceil(float64(quantity) / float64(div))

	return &ProductionResult{
		Consumption: consume,
		Production:  produce,
		Hours:       int(hours),
	}, nil
}

func (b Building) ConsumesList() []ItemID {
	var consume = make(map[ItemID]struct{})
	for i, set := range b.Production {
		if i.IsConsumption {
			consume[i.ItemID] = struct{}{}
		}
		for _, s := range set {
			if s.IsConsumption {
				consume[s.ItemID] = struct{}{}
			}
		}
	}
	var consumeList []ItemID
	for i := range consume {
		consumeList = append(consumeList, i)
	}
	return consumeList
}

func (b Building) ProducesList() []ItemID {
	var produce = make(map[ItemID]struct{})
	for i, set := range b.Production {
		if !i.IsConsumption {
			produce[i.ItemID] = struct{}{}
		}
		for _, s := range set {
			if !s.IsConsumption {
				produce[s.ItemID] = struct{}{}
			}
		}
	}
	var produceList []ItemID
	for i := range produce {
		produceList = append(produceList, i)
	}
	return produceList
}

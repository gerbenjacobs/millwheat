package game

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate stringer -linecomment -trimprefix=Building -type=BuildingType
const (
	BuildingWarehouse BuildingType = iota
	BuildingFarm
	BuildingMill
	BuildingBakery
	BuildingPigFarm // Pig Farm
	BuildingButcher
	BuildingWeaponSmith // Weapon Smith
	BuildingForestry
	BuildingQuarry
	BuildingSawMill // Saw Mill
	BuildingTannery
	BuildingCoalMine // Coal Mine
	BuildingIronMine // Iron Mine
	BuildingBlacksmith
	BuildingArmourSmith // Armour Smith
	BuildingStables
	BuildingVineyard
)

type Buildings map[BuildingType]Building
type BuildingType int

// TownBuilding represents an instance of a building, located in a town
type TownBuilding struct {
	ID                uuid.UUID
	Type              BuildingType
	CurrentLevel      int
	LastCollection    time.Time
	CurrentProduction int
	CreatedAt         time.Time
}

// Building contains all the data for buildings in the game
type Building struct {
	Type        BuildingType
	Name        string
	Description string
	Image       string
	Production  map[ItemSet]ItemSetSlice
	IsGenerator bool
	Mechanics   []BuildingMechanic
	BuildCosts  map[int]BuildingCost
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

// ItemSetSlice is a container for itemset slices used mainly for String() functions
type ItemSetSlice []ItemSet

// ProductionResult contains calculated values for producing an item
type ProductionResult struct {
	Consumption ItemSetSlice
	Production  ItemSetSlice
	Hours       int
}

func (bl Buildings) Sorted() []Building {
	var buildings = make([]Building, len(bl))
	for k, b := range bl {
		b.Type = k
		buildings[k] = b
	}
	sort.SliceStable(buildings, func(i, j int) bool {
		return buildings[i].Name < buildings[j].Name
	})
	return buildings
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
	var consume ItemSetSlice
	var produce ItemSetSlice
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
		maxEfficiency := b.MaxEfficiency(c.ItemID, level)
		maxConsumption := b.MaxConsumption(c.ItemID, level)

		if maxEfficiency == 0 && maxConsumption == 0 {
			consume[i].Quantity = quantity
		} else if maxEfficiency > 0 {
			// efficiency is applied, so consume less; divide
			consume[i].Quantity = int(math.Ceil(float64(quantity) / float64(maxEfficiency)))
		} else {
			// efficiency has priority, then consumption is checked
			if isConsumable {
				consume[i].Quantity = maxConsumption
			} else {
				consume[i].Quantity = maxConsumption * quantity
			}
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

// GeneratedProduct picks the itemID that's produced by a generator building
func (b Building) GeneratedProduct() (ItemID, error) {
	if !b.IsGenerator {
		return "", errors.New("building is not a generator")
	}
	for item, subItems := range b.Production {
		if !item.IsConsumption {
			return item.ItemID, nil
		}

		for _, subItem := range subItems {
			if !subItem.IsConsumption {
				return subItem.ItemID, nil
			}
		}
	}

	return "", errors.New("no generated product found")
}

func (b Building) ActionURL() string {
	if b.IsGenerator {
		return "collect"
	}
	return "produce"
}

func (tb TownBuilding) GetCurrentProduction(b Building) (*ItemSet, error) {
	t := time.Since(tb.LastCollection)
	hours := int(math.Floor(t.Hours()))
	if hours < 1 {
		return nil, errors.New("building not ready for collection")
	}

	if !b.IsGenerator {
		logrus.Warnf("tried to collect from building that is not a generator: %s %s", tb.ID, tb.Type)
		return nil, errors.New("building is not a generator that can be collected from")
	}

	itemID, err := b.GeneratedProduct()
	if err != nil {
		logrus.Errorf("building %s has no generated product", tb.Type)
		return nil, errors.New("building has no generated product")
	}

	quantity := hours * b.MaxProduction(itemID, tb.CurrentLevel)
	if quantity < 1 {
		return nil, errors.New("no produce ready for collection")
	}
	return &ItemSet{
		ItemID:   itemID,
		Quantity: quantity,
	}, nil
}

func (tb TownBuilding) LastCollectionAt() string {
	return tb.LastCollection.Format("2006-01-02 15:04")
}

func (tb TownBuilding) IsWarehouse() bool {
	return tb.Type == BuildingWarehouse
}

func CreateBuilding(building Building, level int) (*ProductionResult, error) {
	costs, ok := building.BuildCosts[level]
	if !ok {
		return nil, errors.New("no build costs found for this level")
	}

	return &ProductionResult{
		Consumption: ItemSetSlice{
			{ItemID: "plank", Quantity: costs.Planks},
			{ItemID: "stone", Quantity: costs.Stones},
		},
		Hours: 1,
	}, nil
}

func RecoverBuilding(building Building, level int) (*ProductionResult, error) {
	costs, ok := building.BuildCosts[level]
	if !ok {
		return nil, errors.New("no build costs found for this level")
	}

	// TODO think about this set up, do we want to return floored value of division by 2?
	return &ProductionResult{
		Consumption: ItemSetSlice{
			{ItemID: "plank", Quantity: costs.Planks / 2},
			{ItemID: "stone", Quantity: costs.Stones / 2},
		},
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
	sort.Slice(consumeList, func(i, j int) bool {
		return consumeList[i] < consumeList[j]
	})
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
	sort.Slice(produceList, func(i, j int) bool {
		return produceList[i] < produceList[j]
	})
	return produceList
}

func (is ItemSet) String() string {
	return fmt.Sprintf("%dx %s", is.Quantity, is.ItemID)
}

func (iss ItemSetSlice) String() string {
	var s []string
	for _, i := range iss {
		s = append(s, fmt.Sprintf("%s", i))
	}

	return strings.Join(s, ", ")
}

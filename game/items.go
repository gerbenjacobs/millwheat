package game

import (
	"fmt"
	"strings"
)

type ItemID string

type Item struct {
	ID          ItemID
	Name        string
	Description string
	Image       string
}

type Items map[ItemID]Item

// WarehouseItem represents an instance of an item in a warehouse
type WarehouseItem struct {
	ItemID   ItemID
	Quantity int
}

func (id ItemID) AsKey() string {
	return strings.Trim(fmt.Sprintf("%#v", id), `""`)
}

// String is used to get a proper name for ItemIDs
func (id ItemID) String() string {
	switch id {
	case "wheat":
		return "Wheat"
	case "flour":
		return "Flour"
	case "bread":
		return "Bread"
	case "stone":
		return "Stone Blocks"
	case "log":
		return "Logs"
	case "plank":
		return "Planks"
	case "wine":
		return "Wine Barrels"
	case "pig":
		return "Pigs"
	case "hide":
		return "Hides"
	case "meat":
		return "Meat"
	case "leather":
		return "Leather Rolls"
	case "iron":
		return "Iron Ore"
	case "coal":
		return "Coal Ore"
	case "iron_bar":
		return "Iron Bars"
	case "horse":
		return "Horses"
	case "leather_armour":
		return "Leather Armour"
	case "iron_platearmour":
		return "Iron Plate Armour"
	case "wooden_shield":
		return "Wooden Shield"
	case "sword":
		return "Sword"
	case "lance":
		return "Lance"
	case "crossbow":
		return "Crossbow"
	default:
		return "Unknown item"
	}
}

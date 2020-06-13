package data

import (
	"github.com/gerbenjacobs/millwheat/game"
)

// WarehouseOrder determines the way the warehouse will be displayed.
var WarehouseOrder = []game.ItemID{
	"stone", "plank", "log", "wheat", "flour", "bread", "wine",
	"pig", "meat", "hide", "leather", "iron", "coal", "iron_bar",
	"woodenshield", "leather_armour", "horse", "iron_platearmour", "sword", "crossbow", "lance",
}

// WarehouseOrderBreakpoints determines when the warehouse starts a new column
var WarehouseOrderBreakpoints = []game.ItemID{"wine", "iron_bar"}

var Items = game.Items{
	"wheat": game.Item{
		ID:          "wheat",
		Name:        "Wheat",
		Description: "A bundle of wheat",
		Image:       "/images/items/wheat.png",
	},
	"flour": game.Item{
		ID:          "flour",
		Name:        "Flour",
		Description: "Coarsely ground grain from a windmill",
		Image:       "/images/items/flour.png",
	},
	"bread": game.Item{
		ID:          "bread",
		Name:        "Bread",
		Description: "A loaf of bread, freshly baked",
		Image:       "/images/items/bread.png",
	},
	"log": game.Item{
		ID:          "log",
		Name:        "Logs",
		Description: "Logs chopped down in forests",
		Image:       "/images/items/log.png",
	},
	"plank": game.Item{
		ID:          "plank",
		Name:        "Planks",
		Description: "Planks sawed from big logs",
		Image:       "/images/items/plank.png",
	},
	"stone": game.Item{
		ID:          "stone",
		Name:        "Stone Blocks",
		Description: "-",
		Image:       "/images/items/stone.png",
	},
	"wine": game.Item{
		ID:          "wine",
		Name:        "Wine barrels",
		Description: "-",
		Image:       "/images/items/wine.png",
	},
	"pig": game.Item{
		ID:          "pig",
		Name:        "Pigs",
		Description: "-",
		Image:       "/images/items/pig.png",
	},
	"meat": game.Item{
		ID:          "meat",
		Name:        "Meat",
		Description: "-",
		Image:       "/images/items/meat.png",
	},
	"hide": game.Item{
		ID:          "hide",
		Name:        "Hides",
		Description: "-",
		Image:       "/images/items/hide.png",
	},
	"leather": game.Item{
		ID:          "leather",
		Name:        "Leather rolls",
		Description: "-",
		Image:       "/images/items/leather.png",
	},
	"iron": game.Item{
		ID:          "iron",
		Name:        "Iron",
		Description: "-",
		Image:       "/images/items/iron.png",
	},
	"coal": game.Item{
		ID:          "coal",
		Name:        "Coal",
		Description: "-",
		Image:       "/images/items/coal.png",
	},
	"iron_bar": game.Item{
		ID:          "iron_bar",
		Name:        "Iron Bars",
		Description: "-",
		Image:       "/images/items/iron_bar.png",
	},
	"horse": game.Item{
		ID:          "horse",
		Name:        "Horse",
		Description: "-",
		Image:       "/images/items/horses.gif",
	},
	"leather_armour": game.Item{
		ID:          "leather_armour",
		Name:        "Leather Armour",
		Description: "-",
		Image:       "/images/items/leather_armour.gif",
	},
	"woodenshield": game.Item{
		ID:          "woodenshield",
		Name:        "Wooden Shield",
		Description: "-",
		Image:       "/images/items/woodenshield.png",
	},
	"iron_platearmour": game.Item{
		ID:          "iron_platearmour",
		Name:        "Iron Plate Armour",
		Description: "Perfect fit for a Knight",
		Image:       "/images/items/iron_armour.gif",
	},
	"sword": game.Item{
		ID:          "sword",
		Name:        "Sword",
		Description: "-",
		Image:       "/images/items/sword.png",
	},
	"crossbow": game.Item{
		ID:          "crossbow",
		Name:        "Crossbow",
		Description: "-",
		Image:       "/images/items/crossbow.gif",
	},
	"lance": game.Item{
		ID:          "lance",
		Name:        "Lance",
		Description: "-",
		Image:       "/images/items/lance.gif",
	},
}

package data

import (
	"github.com/gerbenjacobs/millwheat/game"
)

var Buildings = game.Buildings{
	game.BuildingFarm: {
		Name:        "Farm",
		Description: "Grows wheat in the fields.",
		Image:       "/images/buildings/farm.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "wheat"}: {},
		},
		IsGenerator: true,
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Wheat per hour",
				ItemID: "wheat",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingMill: {
		Name:        "Mill",
		Description: "Mills wheat into bags of flour.",
		Image:       "/images/buildings/mill.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "flour"}: {{ItemID: "wheat", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Wheat per hour",
				ItemID: "wheat",
				Levels: map[int]int{
					1:  1,
					2:  2,
					3:  3,
					4:  4,
					5:  4,
					6:  5,
					7:  6,
					8:  6,
					9:  7,
					10: 8,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Flour per wheat",
				ItemID: "wheat",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Flour per hour",
				ItemID: "flour",
				Levels: map[int]int{
					1:  1,
					2:  2,
					3:  3,
					4:  4,
					5:  8,
					6:  10,
					7:  12,
					8:  18,
					9:  21,
					10: 24,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingBakery: {
		Name:        "Bakery",
		Description: "Bakes bread for the soldiers using flour from the mill.",
		Image:       "/images/buildings/bakery.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "bread"}: {{ItemID: "flour", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Flour per hour",
				ItemID: "flour",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Bread per flour",
				ItemID: "bread",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Bread per hour",
				ItemID: "bread",
				Levels: map[int]int{
					1: 1,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingPigFarm: {
		Name:        "Pig Farm",
		Description: "Raises pigs from piglets with love and a lot of wheat!",
		Image:       "/images/buildings/pigfarm.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "pig"}: {{ItemID: "wheat", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Pigs per hour",
				ItemID: "pig",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
					4: 4,
					5: 4,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Pigs per wheat",
				ItemID: "wheat",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingButcher: {
		Name:        "Butcher",
		Description: "Turns pigs into meat and hide.",
		Image:       "/images/buildings/butcher.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "pig", IsConsumption: true}: {{ItemID: "hide"}, {ItemID: "meat"}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Pigs per hour",
				ItemID: "pig",
				Levels: map[int]int{
					1: 1,
					2: 1,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Hide per pig",
				ItemID: "hide",
				Levels: map[int]int{
					1: 1,
					2: 1,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Meat per pig",
				ItemID: "meat",
				Levels: map[int]int{
					1: 1,
					2: 2,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Hide per hour",
				ItemID: "hide",
				Levels: map[int]int{
					1: 1,
					2: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Meat per hour",
				ItemID: "meat",
				Levels: map[int]int{
					1: 1,
					2: 2,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingWeaponSmith: {
		Name:        "Weapon Smith",
		Description: "Use iron bars and planks to create weaponry.",
		Image:       "/images/buildings/weaponsmith.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "sword"}: {
				{ItemID: "iron_bar", IsConsumption: true},
			},
			{ItemID: "crossbow"}: {
				{ItemID: "plank", IsConsumption: true},
			},
			{ItemID: "lance"}: {
				{ItemID: "iron_bar", IsConsumption: true},
				{ItemID: "plank", IsConsumption: true},
			},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Sword per hour",
				ItemID: "sword",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Crossbow per hour",
				ItemID: "crossbow",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Lance per hour",
				ItemID: "lance",
				Levels: map[int]int{
					1: 1,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingForestry: {
		Name:        "Forestry",
		Description: "Nourishes the forests with saplings and takes out old wood.",
		Image:       "/images/buildings/woodcutter.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "log"}: {},
		},
		IsGenerator: true,
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Logs per hour",
				ItemID: "log",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingSawMill: {
		Name:        "Saw Mill",
		Description: "Saws large logs into planks on a big table saw.",
		Image:       "/images/buildings/sawmill.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "plank"}: {{ItemID: "log", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Logs per hour",
				ItemID: "log",
				Levels: map[int]int{
					1:  1,
					2:  2,
					3:  3,
					4:  4,
					5:  4,
					6:  5,
					7:  6,
					8:  6,
					9:  7,
					10: 8,
				},
			},
			{
				Type:   game.MechanicEfficiency,
				Name:   "Planks per log",
				ItemID: "log",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Planks per hour",
				ItemID: "plank",
				Levels: map[int]int{
					1:  1,
					2:  2,
					3:  3,
					4:  4,
					5:  8,
					6:  10,
					7:  12,
					8:  18,
					9:  21,
					10: 24,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingQuarry: {
		Name:        "Stone Quarry",
		Description: "Quarries the mine for raw stone and turns it into stone blocks.",
		Image:       "/images/buildings/quarry.png",
		IsGenerator: true,
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "stone"}: {},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Stone per hour",
				ItemID: "stone",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingTannery: {
		Name:        "Tannery",
		Description: "Tanner prepares hides for leather production.",
		Image:       "/images/buildings/tannery.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "leather"}: {{ItemID: "hide", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Leather per hour",
				ItemID: "leather",
				Levels: map[int]int{
					1:  1,
					2:  2,
					3:  4,
					4:  6,
					5:  8,
					6:  10,
					7:  12,
					8:  12,
					9:  12,
					10: 12,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingCoalMine: {
		Name:        "Coal Mine",
		Description: "Mines coal from under the ground.",
		Image:       "/images/buildings/coalmine.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "coal"}: {},
		},
		IsGenerator: true,
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Coal per hour",
				ItemID: "coal",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingIronMine: {
		Name:        "Iron Mine",
		Description: "Mines iron from deep into the mountains.",
		Image:       "/images/buildings/ironmine.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "iron"}: {},
		},
		IsGenerator: true,
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Iron per hour",
				ItemID: "iron",
				Levels: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingBlacksmith: {
		Name:        "Blacksmith",
		Description: "Smiths iron bars out of coal and iron ore.",
		Image:       "/images/buildings/blacksmith.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "iron_bar"}: {{ItemID: "coal", IsConsumption: true}, {ItemID: "iron", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Coal per iron bar",
				ItemID: "coal",
				Levels: map[int]int{
					1: 4,
					2: 3,
					3: 2,
					4: 1,
					5: 1,
				},
			},
			{
				Type:   game.MechanicConsumption,
				Name:   "Iron per iron bar",
				ItemID: "iron",
				Levels: map[int]int{
					1: 2,
					2: 2,
					3: 2,
					4: 1,
					5: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Iron Bars per hour",
				ItemID: "iron_bar",
				Levels: map[int]int{
					1: 1,
					2: 1,
					3: 2,
					4: 2,
					5: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {3, 6},
			3: {5, 9},
			4: {8, 12},
			5: {13, 15},
		},
	},
	game.BuildingArmourSmith: {
		Name:        "Armour Smith",
		Description: "Use leather, plank and iron bars to create crude armour.",
		Image:       "https://www.knightsandmerchants.net/application/files/1015/6823/6439/armoryworkshop.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "wooden_shield"}: {
				{ItemID: "plank", IsConsumption: true},
			},
			{ItemID: "leather_armour"}: {
				{ItemID: "leather", IsConsumption: true},
			},
			{ItemID: "iron_platearmour"}: {
				{ItemID: "iron_bar", IsConsumption: true},
			},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicOutput,
				Name:   "Wooden Shields per hour",
				ItemID: "wooden_shield",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Leather Armour per hour",
				ItemID: "leather_armour",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Iron Plate Armour per hour",
				ItemID: "iron_platearmour",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  1,
					4:  1,
					5:  2,
					6:  2,
					7:  2,
					8:  3,
					9:  3,
					10: 3,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingStables: {
		Name:        "Stables",
		Description: "Horses bred for knights. Love a good dose of wheat!",
		Image:       "https://www.knightsandmerchants.net/application/files/7715/6823/6448/stables.png",
		Production: map[game.ItemSet]game.ItemSetSlice{
			{ItemID: "horse"}: {{ItemID: "wheat", IsConsumption: true}},
		},
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicConsumption,
				Name:   "Wheat per horse",
				ItemID: "wheat",
				Levels: map[int]int{
					1:  4,
					2:  4,
					3:  3,
					4:  3,
					5:  3,
					6:  2,
					7:  2,
					8:  1,
					9:  1,
					10: 1,
				},
			},
			{
				Type:   game.MechanicOutput,
				Name:   "Horse per hour",
				ItemID: "horse",
				Levels: map[int]int{
					1:  1,
					2:  1,
					3:  2,
					4:  2,
					5:  3,
					6:  3,
					7:  4,
					8:  4,
					9:  5,
					10: 5,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 18},
			4: {5, 50},
			5: {7, 75},
		},
	},
	game.BuildingWarehouse: {
		Name:        "Warehouse",
		Description: "Giant warehouse shared by the guilds of the town.",
		Image:       "https://www.knightsandmerchants.net/application/files/3515/6823/6449/storehouse.png",
		Mechanics: []game.BuildingMechanic{
			{
				Type:   game.MechanicEfficiency,
				Name:   "Max quantity",
				ItemID: "slots",
				Levels: map[int]int{
					1: 100,
					2: 130,
					3: 175,
					4: 225,
					5: 300,
				},
			},
		},
		BuildCosts: map[int]game.BuildingCost{
			1: {1, 3},
			2: {2, 6},
			3: {3, 15},
			4: {5, 50},
			5: {7, 75},
		},
	},
}

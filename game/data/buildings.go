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
		Description: "Nourishes the forests with saplings and takes out old wood",
		Image:       "https://www.knightsandmerchants.net/application/files/7315/6823/6438/woodcutters.png",
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
		Description: "Saws large logs into planks on a big table saw",
		Image:       "https://www.knightsandmerchants.net/application/files/4715/6823/6447/sawmill.png",
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
		Image:       "https://www.knightsandmerchants.net/application/files/7115/6823/6446/quarry.png",
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
		Image:       "https://www.knightsandmerchants.net/application/files/3615/6823/6450/tannery.png",
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
}

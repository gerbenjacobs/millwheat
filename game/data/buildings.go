package data

import (
	"time"

	"github.com/gerbenjacobs/millwheat/game"
)

var Buildings = game.Buildings{
	game.BuildingFarm: {
		Name:        "Farm",
		Description: "Grows wheat in the fields.",
		Image:       "/images/buildings/farm.png",
		Production: map[game.ItemSet][]game.ItemSet{
			{ItemID: "wheat"}: {},
		},
		IsGenerator:    true,
		LastCollection: time.Now().Add(-2 * time.Hour).UTC(),
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
		Production: map[game.ItemSet][]game.ItemSet{
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
		Production: map[game.ItemSet][]game.ItemSet{
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
		Production: map[game.ItemSet][]game.ItemSet{
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
		Production: map[game.ItemSet][]game.ItemSet{
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
		Production: map[game.ItemSet][]game.ItemSet{
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
}

package data

import (
	"github.com/gerbenjacobs/millwheat/game"
)

var Buildings = game.Buildings{
	game.BuildingFarm: {
		Name:        "Farm",
		Description: "Grows wheat in the fields.",
		Image:       "https://www.knightsandmerchants.net/application/files/7515/6823/6441/farm.png",
		Consumes:    nil,
		Produces:    []game.ItemID{"wheat"},
		IsGenerator: true,
		Mechanics: []game.BuildingMechanic{
			{
				Type: game.MechanicOutput,
				Name: "Wheat per hour",
				Levels: map[int]int{
					1: 1,
					2: 2,
				},
			},
		},
		BuildCosts: map[int]game.BuildCost{
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
		Image:       "https://www.knightsandmerchants.net/application/files/9415/6823/6446/mill.png",
		Consumes:    []game.ItemID{"wheat"},
		Produces:    []game.ItemID{"flour"},
		Mechanics: []game.BuildingMechanic{
			{
				Type: game.MechanicConsumption,
				Name: "Wheat per hour",
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
				Type: game.MechanicEfficiency,
				Name: "Flour per wheat",
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
				Type: game.MechanicOutput,
				Name: "Flour per hour",
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
	},
	game.BuildingBakery: {
		Name:        "Bakery",
		Description: "Bakes bread for the soldiers using flour from the mill.",
		Image:       "https://www.knightsandmerchants.net/application/files/1215/6823/6439/bakery.png",
		Consumes:    []game.ItemID{"flour"},
		Produces:    []game.ItemID{"bread"},
		Mechanics: []game.BuildingMechanic{
			{
				Type: game.MechanicConsumption,
				Name: "Flour per hour",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type: game.MechanicEfficiency,
				Name: "Bread per flour",
				Levels: map[int]int{
					1: 1,
				},
			},
			{
				Type: game.MechanicOutput,
				Name: "Bread per hour",
				Levels: map[int]int{
					1: 1,
				},
			},
		},
	},
}

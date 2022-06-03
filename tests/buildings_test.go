// CreateProduct tests are in a separate folder to prevent import cycles with game/data
package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

func TestBuilding_CreateProduct(t *testing.T) {
	type productRequest struct {
		product  game.ItemID
		quantity int
		level    int
	}
	tests := []struct {
		building game.Building
		req      productRequest
		want     *game.ProductionResult
		wantErr  bool
	}{
		{
			building: gamedata.Buildings[game.BuildingMill],
			req:      productRequest{"flour", 1, 1},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "wheat", Quantity: 1, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "flour", Quantity: 1},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingMill],
			req:      productRequest{"flour", 2, 2},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "wheat", Quantity: 2, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "flour", Quantity: 2},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingMill],
			req:      productRequest{"flour", 8, 5},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "wheat", Quantity: 4, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "flour", Quantity: 8},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingMill],
			req:      productRequest{"flour", 9, 5},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "wheat", Quantity: 5, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "flour", Quantity: 9},
				},
				Hours: 2,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 1, 1},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 1, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 1},
					{ItemID: "meat", Quantity: 1},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 1, 2},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 1, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 1},
					{ItemID: "meat", Quantity: 2},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 1, 3},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 1, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 2},
					{ItemID: "meat", Quantity: 2},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 2, 4},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 2, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 8},
					{ItemID: "meat", Quantity: 10},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 2, 5},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 2, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 10},
					{ItemID: "meat", Quantity: 10},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingButcher],
			req:      productRequest{"pig", 6, 1},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "pig", Quantity: 6, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "hide", Quantity: 6},
					{ItemID: "meat", Quantity: 6},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingWeaponSmith],
			req:      productRequest{"lance", 1, 1},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "iron_bar", Quantity: 1, IsConsumption: true},
					{ItemID: "plank", Quantity: 1, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "lance", Quantity: 1},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingBlacksmith],
			req:      productRequest{"iron_bar", 1, 1},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "coal", Quantity: 4, IsConsumption: true},
					{ItemID: "iron", Quantity: 2, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "iron_bar", Quantity: 1},
				},
				Hours: 1,
			},
		},
		{
			building: gamedata.Buildings[game.BuildingBlacksmith],
			req:      productRequest{"iron_bar", 2, 3},
			want: &game.ProductionResult{
				Consumption: []game.ItemSet{
					{ItemID: "coal", Quantity: 4, IsConsumption: true},
					{ItemID: "iron", Quantity: 4, IsConsumption: true},
				},
				Production: []game.ItemSet{
					{ItemID: "iron_bar", Quantity: 2},
				},
				Hours: 1,
			},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%d %s at level %d %s", tt.req.quantity, tt.req.product, tt.req.level, tt.building.Name)
		t.Run(name, func(t *testing.T) {
			got, err := tt.building.CreateProduct(tt.req.product, tt.req.quantity, tt.req.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateProduct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package tests

import (
	"reflect"
	"testing"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

func TestBuilding_CreateProduct(t *testing.T) {

	type args struct {
		product  game.ItemID
		quantity int
		level    int
	}
	tests := []struct {
		name     string
		building game.Building
		args     args
		want     *game.ProductionResult
		wantErr  bool
	}{
		{
			name:     "1 wheat into flour in level 1 mill",
			building: gamedata.Buildings[game.BuildingMill],
			args:     args{"flour", 1, 1},
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
			name:     "2 wheat into flour in level 2 mill",
			building: gamedata.Buildings[game.BuildingMill],
			args:     args{"flour", 2, 2},
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
			name:     "8 wheat into flour in level 5 mill",
			building: gamedata.Buildings[game.BuildingMill],
			args:     args{"flour", 8, 5},
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
			name:     "9 wheat into flour in level 5 mill",
			building: gamedata.Buildings[game.BuildingMill],
			args:     args{"flour", 9, 5},
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
			name:     "1 pig for level 1 butcher",
			building: gamedata.Buildings[game.BuildingButcher],
			args:     args{"pig", 1, 1},
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
			name:     "1 pig for level 2 butcher",
			building: gamedata.Buildings[game.BuildingButcher],
			args:     args{"pig", 1, 2},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.building.CreateProduct(tt.args.product, tt.args.quantity, tt.args.level)
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

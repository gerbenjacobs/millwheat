package tests

import (
	"testing"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

func TestBuilding_MaxProduction(t *testing.T) {
	type args struct {
		itemID game.ItemID
		level  int
	}
	tests := []struct {
		name     string
		building game.Building
		args     args
		want     int
	}{
		{
			name:     "Max production for lvl 1 Butcher",
			building: gamedata.Buildings[game.BuildingButcher],
			args:     args{itemID: "pig", level: 1},
			want:     1,
		},
		{
			name:     "Max production for lvl 4 Butcher",
			building: gamedata.Buildings[game.BuildingButcher],
			args:     args{itemID: "pig", level: 4},
			want:     2,
		},
		{
			name:     "Max production for Crossbow at lvl 1 Weapon Smith",
			building: gamedata.Buildings[game.BuildingWeaponSmith],
			args:     args{itemID: "crossbow", level: 1},
			want:     1,
		},
		{
			name:     "Max production for Crossbow at lvl 2 Weapon Smith",
			building: gamedata.Buildings[game.BuildingWeaponSmith],
			args:     args{itemID: "crossbow", level: 2},
			want:     2,
		},
		{
			name:     "Max production at lvl 1 Pig Farm",
			building: gamedata.Buildings[game.BuildingPigFarm],
			args:     args{itemID: "pig", level: 1},
			want:     1,
		},
		{
			name:     "Max production at lvl 2 Pig Farm",
			building: gamedata.Buildings[game.BuildingPigFarm],
			args:     args{itemID: "pig", level: 2},
			want:     2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.building.MaxProduction(tt.args.itemID, tt.args.level); got != tt.want {
				t.Errorf("MaxProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

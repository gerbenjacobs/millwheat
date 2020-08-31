package game

import (
	"reflect"
	"testing"
)

func TestCalculateWarriorCosts(t *testing.T) {
	type args struct {
		warriorType WarriorType
		quantity    int
	}
	tests := []struct {
		name    string
		args    args
		want    ItemSetSlice
		wantErr bool
	}{
		{
			name: "create 1 sword warrior",
			args: args{warriorType: WarriorSword, quantity: 1},
			want: ItemSetSlice{
				{ItemID: "sword", Quantity: 1},
				{ItemID: "wooden_shield", Quantity: 1},
				{ItemID: "wine", Quantity: 1},
				{ItemID: "bread", Quantity: 1},
			},
			wantErr: false,
		},
		{
			name: "create 2 sword warrior",
			args: args{warriorType: WarriorSword, quantity: 2},
			want: ItemSetSlice{
				{ItemID: "sword", Quantity: 2},
				{ItemID: "wooden_shield", Quantity: 2},
				{ItemID: "wine", Quantity: 2},
				{ItemID: "bread", Quantity: 2},
			},
			wantErr: false,
		},
		{
			name: "create 1 crossbow warrior",
			args: args{warriorType: WarriorCrossbow, quantity: 1},
			want: ItemSetSlice{
				{ItemID: "crossbow", Quantity: 1},
				{ItemID: "leather_armour", Quantity: 1},
				{ItemID: "wine", Quantity: 1},
				{ItemID: "bread", Quantity: 1},
			},
			wantErr: false,
		},
		{
			name: "create 1 lance warrior",
			args: args{warriorType: WarriorLance, quantity: 1},
			want: ItemSetSlice{
				{ItemID: "lance", Quantity: 1},
				{ItemID: "iron_platearmour", Quantity: 1},
				{ItemID: "horse", Quantity: 1},
				{ItemID: "wine", Quantity: 1},
				{ItemID: "meat", Quantity: 1},
			},
			wantErr: false,
		},
		{
			name: "create 14 lance warrior",
			args: args{warriorType: WarriorLance, quantity: 14},
			want: ItemSetSlice{
				{ItemID: "lance", Quantity: 14},
				{ItemID: "iron_platearmour", Quantity: 14},
				{ItemID: "horse", Quantity: 14},
				{ItemID: "wine", Quantity: 14},
				{ItemID: "meat", Quantity: 14},
			},
			wantErr: false,
		},
		{
			name:    "create 1 wrong warrior type",
			args:    args{warriorType: 88, quantity: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateWarriorCosts(tt.args.warriorType, tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateWarriorCosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateWarriorCosts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

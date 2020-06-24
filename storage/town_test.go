package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/gerbenjacobs/millwheat/game"
)

func TestTownRepository_TakeFromWarehouse(t1 *testing.T) {
	town1 := uuid.New()
	town2 := uuid.New()

	itemsTaken := []game.ItemSet{{ItemID: "wheat", Quantity: 2}}

	wantedItems1 := map[game.ItemID]game.WarehouseItem{
		"wheat": {ItemID: "wheat", Quantity: 8},
	}
	wantedItems2 := map[game.ItemID]game.WarehouseItem{
		"wheat": {ItemID: "wheat", Quantity: 10},
	}

	t := &TownRepository{
		warehouseCache: cache.New(cache.NoExpiration, 0),
	}
	t.warehouseCache.Set(town1.String(), map[game.ItemID]game.WarehouseItem{"wheat": {ItemID: "wheat", Quantity: 10}}, cache.NoExpiration)
	t.warehouseCache.Set(town2.String(), map[game.ItemID]game.WarehouseItem{"wheat": {ItemID: "wheat", Quantity: 10}}, cache.NoExpiration)

	if err := t.TakeFromWarehouse(context.Background(), town1, itemsTaken); err != nil {
		t1.Errorf("TakeFromWarehouse() error = %v", err)
	}

	// check town 1
	items, err := t.WarehouseItems(context.Background(), town1)
	if err != nil {
		t1.Errorf("WarehouseItems() error = %v", err)
	}

	if !reflect.DeepEqual(items, wantedItems1) {
		t1.Errorf("no match between items got: %v  want: %v", items, wantedItems1)
	}

	// check town 2
	items2, err := t.WarehouseItems(context.Background(), town2)
	if err != nil {
		t1.Errorf("WarehouseItems() error = %v", err)
	}

	if !reflect.DeepEqual(items2, wantedItems2) {
		t1.Errorf("no match between items got: %v  want: %v", items2, wantedItems2)
	}

}

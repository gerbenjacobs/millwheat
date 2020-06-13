package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type TownRepository struct {
	towns      map[uuid.UUID]*game.Town
	warehouses map[uuid.UUID]map[game.ItemID]game.WarehouseItem
}

func NewTownRepository(towns map[uuid.UUID]*game.Town) *TownRepository {
	defaultWarehouse := map[game.ItemID]game.WarehouseItem{
		game.ItemID("wheat"): {ItemID: game.ItemID("wheat"), Quantity: 10},
		game.ItemID("flour"): {ItemID: game.ItemID("flour"), Quantity: 10},
	}
	warehouses := make(map[uuid.UUID]map[game.ItemID]game.WarehouseItem)
	for id := range towns {
		warehouses[id] = defaultWarehouse
	}
	return &TownRepository{towns: towns, warehouses: warehouses}
}

func (t *TownRepository) Get(_ context.Context, id uuid.UUID) (*game.Town, error) {
	if town, ok := t.towns[id]; ok {
		return town, nil
	}

	return nil, errors.New("town not found")
}

func (t *TownRepository) WarehouseItems(_ context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error) {
	if wh, ok := t.warehouses[townID]; ok {
		return wh, nil
	}

	return nil, errors.New("warehouse not found")
}

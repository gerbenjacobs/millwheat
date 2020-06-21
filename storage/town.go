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
		game.ItemID("stone"):    {ItemID: game.ItemID("stone"), Quantity: 5},
		game.ItemID("plank"):    {ItemID: game.ItemID("plank"), Quantity: 12},
		game.ItemID("wheat"):    {ItemID: game.ItemID("wheat"), Quantity: 10},
		game.ItemID("flour"):    {ItemID: game.ItemID("flour"), Quantity: 4},
		game.ItemID("iron_bar"): {ItemID: game.ItemID("iron_bar"), Quantity: 4},
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

func (t *TownRepository) ItemsInWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) bool {
	for _, is := range items {
		i, ok := t.warehouses[townID][is.ItemID]
		if !ok {
			// if item not found
			return false
		}
		if i.Quantity < is.Quantity {
			// if not enough quantity for this item
			return false
		}
	}

	return true
}

func (t *TownRepository) TakeFromWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) error {
	for _, is := range items {
		i, ok := t.warehouses[townID][is.ItemID]
		if !ok {
			// if item not found
			return errors.New("item not found")
		}
		if i.Quantity < is.Quantity {
			// if not enough quantity for this item
			return errors.New("not enough quantity")
		}

		t.warehouses[townID][is.ItemID] = game.WarehouseItem{
			ItemID:   i.ItemID,
			Quantity: i.Quantity - is.Quantity,
		}
	}

	return nil
}

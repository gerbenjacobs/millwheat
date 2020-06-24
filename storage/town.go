package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/gerbenjacobs/millwheat/game"
)

type TownRepository struct {
	towns          map[uuid.UUID]*game.Town
	warehouseCache *cache.Cache
}

// defaultWarehouse returns a new map, to prevent pointer issues
func defaultWarehouse() map[game.ItemID]game.WarehouseItem {
	return map[game.ItemID]game.WarehouseItem{
		"stone":    {ItemID: "stone", Quantity: 5},
		"plank":    {ItemID: "plank", Quantity: 12},
		"wheat":    {ItemID: "wheat", Quantity: 10},
		"flour":    {ItemID: "flour", Quantity: 4},
		"iron_bar": {ItemID: "iron_bar", Quantity: 4},
	}
}

func NewTownRepository(towns map[uuid.UUID]*game.Town) *TownRepository {
	c := cache.New(cache.NoExpiration, 0)
	for id := range towns {
		c.Set(id.String(), defaultWarehouse(), cache.NoExpiration)
	}
	return &TownRepository{towns: towns, warehouseCache: c}
}

func (t *TownRepository) Get(_ context.Context, id uuid.UUID) (*game.Town, error) {
	if town, ok := t.towns[id]; ok {
		return town, nil
	}

	return nil, errors.New("town not found")
}

func (t *TownRepository) WarehouseItems(_ context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error) {
	wh, ok := t.warehouseCache.Get(townID.String())
	if !ok {
		return nil, errors.New("warehouse not found")
	}

	return wh.(map[game.ItemID]game.WarehouseItem), nil
}

func (t *TownRepository) ItemsInWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) bool {
	wh, err := t.WarehouseItems(ctx, townID)
	if err != nil {
		return false
	}
	for _, is := range items {
		i, ok := wh[is.ItemID]
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
	wh, err := t.WarehouseItems(ctx, townID)
	if err != nil {
		return err
	}

	for _, is := range items {
		i, ok := wh[is.ItemID]
		if !ok {
			// if item not found
			return errors.New("item not found")
		}
		if i.Quantity < is.Quantity {
			// if not enough quantity for this item
			return errors.New("not enough quantity")
		}

		wh[is.ItemID] = game.WarehouseItem{
			ItemID:   i.ItemID,
			Quantity: i.Quantity - is.Quantity,
		}
	}

	t.warehouseCache.Set(townID.String(), wh, cache.NoExpiration)
	return nil
}

func (t *TownRepository) GiveToWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) error {
	wh, err := t.WarehouseItems(ctx, townID)
	if err != nil {
		return err
	}
	for _, is := range items {
		i, ok := wh[is.ItemID]
		if ok && i.Quantity+is.Quantity > 100 { // TODO fix hardcode warehouse upper limit
			// set quantity to upper limit
			is.Quantity = 100 - i.Quantity
		}

		wh[is.ItemID] = game.WarehouseItem{
			ItemID:   i.ItemID,
			Quantity: i.Quantity + is.Quantity,
		}
	}

	t.warehouseCache.Set(townID.String(), wh, cache.NoExpiration)
	return nil
}

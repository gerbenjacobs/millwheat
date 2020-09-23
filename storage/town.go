package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/game/data"
)

var (
	CacheDurationTown = 6 * time.Hour

	defaultWarehouseLimit = 100
)

type TownRepository struct {
	db        *sql.DB
	townCache *cache.Cache
}

// defaultWarehouse returns a new map, to prevent pointer issues
func defaultWarehouse() map[game.ItemID]game.WarehouseItem {
	return map[game.ItemID]game.WarehouseItem{
		"stone":    {ItemID: "stone", Quantity: 100},
		"plank":    {ItemID: "plank", Quantity: 100},
		"wheat":    {ItemID: "wheat", Quantity: 10},
		"flour":    {ItemID: "flour", Quantity: 4},
		"iron_bar": {ItemID: "iron_bar", Quantity: 4},
	}
}

func NewTownRepository(db *sql.DB) *TownRepository {
	c := cache.New(CacheDurationTown, time.Hour)

	return &TownRepository{db: db, townCache: c}
}

func (t *TownRepository) Get(ctx context.Context, id uuid.UUID) (town *game.Town, err error) {
	tc, ok := t.townCache.Get(id.String())
	if !ok {
		town, err = t.getTownFromDatabase(ctx, id)
		if err != nil {
			return nil, err
		}
	} else {
		town, _ = tc.(*game.Town)
	}

	// calculate current production for generator buildings
	for id, tb := range town.Buildings {
		b, ok := data.Buildings[tb.Type]
		if !ok || !b.IsGenerator {
			continue
		}
		cp, err := tb.GetCurrentProduction(b)
		if err != nil {
			continue
		}
		tb.CurrentProduction = cp.Quantity
		town.Buildings[id] = tb
	}

	return town, nil
}

func (t *TownRepository) AddBuilding(ctx context.Context, townID uuid.UUID, buildingType game.BuildingType) error {
	tb := game.TownBuilding{
		ID:             uuid.New(),
		Type:           buildingType,
		CurrentLevel:   1,
		LastCollection: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	return t.addBuildingToDatabase(ctx, townID, tb)
}

func (t *TownRepository) doesHaveBuilding(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) (*game.TownBuilding, error) {
	town, err := t.Get(ctx, townID)
	if err != nil {
		return nil, fmt.Errorf("town not found: %w", err)
	}

	for _, b := range town.Buildings {
		if b.ID == buildingID {
			return &b, nil
		}
	}

	return nil, errors.New("building not found")
}

func (t *TownRepository) UpgradeBuilding(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error {
	cb, err := t.doesHaveBuilding(ctx, townID, buildingID)
	if err != nil {
		return err
	}

	cb.CurrentLevel = cb.CurrentLevel + 1
	return t.upgradeBuildingInDatabase(ctx, townID, *cb)
}

func (t *TownRepository) RemoveBuilding(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error {
	_, err := t.doesHaveBuilding(ctx, townID, buildingID)
	if err != nil {
		return err
	}

	return t.removeBuildingInDatabase(ctx, townID, buildingID)
}

func (t *TownRepository) BuildingCollected(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error {
	cb, err := t.doesHaveBuilding(ctx, townID, buildingID)
	if err != nil {
		return err
	}

	cb.LastCollection = time.Now().UTC()
	cb.CurrentProduction = 0
	return t.updateBuildingCollection(ctx, townID, *cb)
}

func (t *TownRepository) WarehouseItems(ctx context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error) {
	town, err := t.Get(ctx, townID)
	if err != nil {
		return nil, err
	}
	return town.Warehouse, nil
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

	// do changes in temporary warehouse struct
	newWh := make(map[game.ItemID]game.WarehouseItem)
	for _, i := range wh {
		newWh[i.ItemID] = i
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

		newWh[i.ItemID] = game.WarehouseItem{
			ItemID:   i.ItemID,
			Quantity: i.Quantity - is.Quantity,
		}
	}

	return t.updateWarehouseInDatabase(ctx, townID, newWh)
}

func (t *TownRepository) GiveToWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) error {
	wh, err := t.WarehouseItems(ctx, townID)
	if err != nil {
		return err
	}

	currentLimit := t.warehouseLimit(townID)
	for _, is := range items {
		i, ok := wh[is.ItemID]
		if !ok {
			i.Quantity = 0
			i.ItemID = is.ItemID
		}
		if i.Quantity+is.Quantity > currentLimit {
			// set quantity to upper limit
			is.Quantity = currentLimit - i.Quantity
		}

		wh[i.ItemID] = game.WarehouseItem{
			ItemID:   i.ItemID,
			Quantity: i.Quantity + is.Quantity,
		}
	}

	return t.updateWarehouseInDatabase(ctx, townID, wh)
}

func (t *TownRepository) warehouseLimit(townID uuid.UUID) int {
	town, err := t.Get(context.Background(), townID)
	if err != nil {
		return defaultWarehouseLimit
	}
	for _, tb := range town.Buildings {
		if tb.IsWarehouse() {
			return data.Buildings[tb.Type].MaxEfficiency("slots", tb.CurrentLevel)
		}
	}

	return defaultWarehouseLimit
}

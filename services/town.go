package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/storage"
)

type TownSvc struct {
	storage storage.TownStorage
}

func NewTownSvc(storage storage.TownStorage) *TownSvc {
	return &TownSvc{storage: storage}
}

func (t *TownSvc) Town(ctx context.Context, id uuid.UUID) (*game.Town, error) {
	return t.storage.Get(ctx, id)
}

func (t *TownSvc) AddBuilding(ctx context.Context, buildingType game.BuildingType) error {
	return t.storage.AddBuilding(ctx, TownFromContext(ctx), buildingType)
}

func (t *TownSvc) UpgradeBuilding(ctx context.Context, buildingID uuid.UUID) error {
	return t.storage.UpgradeBuilding(ctx, TownFromContext(ctx), buildingID)
}

func (t *TownSvc) RemoveBuilding(ctx context.Context, buildingID uuid.UUID) error {
	return t.storage.RemoveBuilding(ctx, TownFromContext(ctx), buildingID)
}

func (t *TownSvc) BuildingCollected(ctx context.Context, buildingID uuid.UUID) error {
	return t.storage.BuildingCollected(ctx, TownFromContext(ctx), buildingID)
}

func (t *TownSvc) Warehouse(ctx context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error) {
	return t.storage.WarehouseItems(ctx, townID)
}

func (t *TownSvc) ItemsInWarehouse(ctx context.Context, items []game.ItemSet) bool {
	return t.storage.ItemsInWarehouse(ctx, TownFromContext(ctx), items)
}

func (t *TownSvc) TakeFromWarehouse(ctx context.Context, items []game.ItemSet) error {
	return t.storage.TakeFromWarehouse(ctx, TownFromContext(ctx), items)
}

func (t *TownSvc) GiveToWarehouse(ctx context.Context, items []game.ItemSet) error {
	return t.storage.GiveToWarehouse(ctx, TownFromContext(ctx), items)
}

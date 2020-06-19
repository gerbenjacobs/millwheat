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

func (t *TownSvc) Warehouse(ctx context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error) {
	return t.storage.WarehouseItems(ctx, townID)
}

func (t *TownSvc) ItemsInWarehouse(ctx context.Context, items []game.ItemSet) bool {
	return t.storage.ItemsInWarehouse(ctx, TownFromContext(ctx), items)
}

func (t *TownSvc) TakeFromWarehouse(ctx context.Context, items []game.ItemSet) error {
	return t.storage.TakeFromWarehouse(ctx, TownFromContext(ctx), items)
}

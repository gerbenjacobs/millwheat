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

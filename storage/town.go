package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type TownRepository struct {
	towns map[uuid.UUID]*game.Town
}

func NewTownRepository(towns map[uuid.UUID]*game.Town) *TownRepository {
	return &TownRepository{towns: towns}
}

func (t *TownRepository) Get(_ context.Context, id uuid.UUID) (*game.Town, error) {
	if town, ok := t.towns[id]; ok {
		return town, nil
	}

	return nil, errors.New("town not found")
}

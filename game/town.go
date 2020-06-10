package game

import (
	"time"

	"github.com/google/uuid"
)

type Towns map[uuid.UUID]*Town

type Town struct {
	ID        uuid.UUID
	Owner     uuid.UUID
	Name      string
	Buildings []TownBuilding
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Town) FormattedCreatedAt() string {
	return t.CreatedAt.Format("2006-01-02 15:04")
}

func (t *Town) FormattedUpdatedAt() string {
	return t.UpdatedAt.Format("2006-01-02 15:04")
}

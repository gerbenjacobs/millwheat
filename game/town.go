package game

import (
	"time"

	"github.com/google/uuid"
)

const (
	BuildingFarm BuildingType = iota
	BuildingMill
	BuildingBakery
)

type BuildingType int

type Town struct {
	ID        uuid.UUID
	Owner     uuid.UUID
	Name      string
	Buildings []Building
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Building struct {
	Type         BuildingType
	CurrentLevel int
}

func (t *Town) FormattedCreatedAt() string {
	return t.CreatedAt.Format("2006-01-02 15:04")
}

func (t *Town) FormattedUpdatedAt() string {
	return t.UpdatedAt.Format("2006-01-02 15:04")
}

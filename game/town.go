package game

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type Towns map[uuid.UUID]*Town

type Town struct {
	ID        uuid.UUID
	Owner     uuid.UUID
	Name      string
	Buildings map[uuid.UUID]TownBuilding
	Warehouse map[ItemID]WarehouseItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Town) FormattedCreatedAt() string {
	return t.CreatedAt.Format("2006-01-02 15:04")
}

func (t *Town) FormattedUpdatedAt() string {
	return t.UpdatedAt.Format("2006-01-02 15:04")
}

func (t *Town) OrderedBuildings() []TownBuilding {
	var buildings []TownBuilding
	for _, b := range t.Buildings {
		buildings = append(buildings, b)
	}

	sort.Slice(buildings, func(i, j int) bool {
		return buildings[i].CreatedAt.Before(buildings[j].CreatedAt)
	})

	return buildings
}

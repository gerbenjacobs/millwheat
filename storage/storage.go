package storage

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/game"
)

type UserStorage interface {
	Create(ctx context.Context, user *app.User) error
	Read(ctx context.Context, userID uuid.UUID) (*app.User, error)
	Login(ctx context.Context, email, password string) (*app.User, error)
	Update(ctx context.Context, user *app.User) (*app.User, error)
}

type TownStorage interface {
	Create(ctx context.Context, owner uuid.UUID, townName string) (*game.Town, error)

	Get(ctx context.Context, id uuid.UUID) (*game.Town, error)
	AddBuilding(ctx context.Context, townID uuid.UUID, buildingType game.BuildingType) error
	UpgradeBuilding(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error
	RemoveBuilding(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error
	BuildingCollected(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error

	WarehouseItems(ctx context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error)
	ItemsInWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) bool
	TakeFromWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) error
	GiveToWarehouse(ctx context.Context, townID uuid.UUID, items []game.ItemSet) error
}

type ProductionStorage interface {
	ProductJobsByTown(ctx context.Context, townID uuid.UUID) map[uuid.UUID][]*game.Job
	QueuedBuildings(ctx context.Context, townID uuid.UUID) []*game.Job
	CreateJob(ctx context.Context, townID uuid.UUID, job *game.Job) error
	UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status game.JobStatus) error
	CancelJob(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) error
	RevertJobResources(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) ([]game.ItemSet, error)

	JobsCompleted(ctx context.Context) map[uuid.UUID][]*game.Job
	ReshuffleQueue(ctx context.Context, townID uuid.UUID)
}

type BattleStorage interface {
	AddWarrior(ctx context.Context, battleId, armyId, townId uuid.UUID, warriorType game.WarriorType, quantity int) error
	WarriorsFromTown(ctx context.Context, townId, battleId uuid.UUID) ([]game.Warrior, error)
	AllWarriorsForBattle(ctx context.Context, battleId uuid.UUID) ([]game.Army, error)
	CurrentWarriors(ctx context.Context, battleId, armyId, townId uuid.UUID) ([]game.Warrior, error)
}

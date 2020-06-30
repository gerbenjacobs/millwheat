package services

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/game"
)

const (
	CtxKeyUserID = iota
	CtxKeyTownID
)

type UserService interface {
	Add(ctx context.Context, user *app.User) error
	User(ctx context.Context, userID uuid.UUID) (*app.User, error)
	Login(ctx context.Context, email, password string) (*app.User, error)
	Update(ctx context.Context, user *app.User) (*app.User, error)
}

type TownService interface {
	Town(ctx context.Context, id uuid.UUID) (*game.Town, error)
	AddBuilding(ctx context.Context, buildingType game.BuildingType) error
	UpgradeBuilding(ctx context.Context, buildingID uuid.UUID) error

	Warehouse(ctx context.Context, townID uuid.UUID) (map[game.ItemID]game.WarehouseItem, error)
	ItemsInWarehouse(ctx context.Context, items []game.ItemSet) bool
	TakeFromWarehouse(ctx context.Context, items []game.ItemSet) error
	GiveToWarehouse(ctx context.Context, items []game.ItemSet) error
}

type ProductionService interface {
	QueuedJobs(ctx context.Context) map[uuid.UUID][]*game.Job
	QueuedBuildings(ctx context.Context) []*game.Job
	CreateJob(ctx context.Context, job *game.InputJob) error
	UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status game.JobStatus) error
	CancelJob(ctx context.Context, jobID uuid.UUID) error
	RevertJobResources(ctx context.Context, jobID uuid.UUID) ([]game.ItemSet, error)

	JobsCompleted(ctx context.Context) map[uuid.UUID][]*game.Job
	ReshuffleQueue(ctx context.Context)
}

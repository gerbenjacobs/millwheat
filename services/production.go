package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/storage"
)

type ProductionSvc struct {
	storage storage.ProductionStorage
}

func NewProductionSvc(storage storage.ProductionStorage) *ProductionSvc {
	return &ProductionSvc{storage: storage}
}

func (p *ProductionSvc) QueuedJobs(ctx context.Context) []*game.Job {
	return p.storage.QueuedJobs(ctx, TownFromContext(ctx))
}

func (p *ProductionSvc) CreateJob(ctx context.Context, job *game.Job) error {
	job.ID = uuid.New()
	job.Created = time.Now().UTC()

	return p.storage.CreateJob(ctx, TownFromContext(ctx), job)
}

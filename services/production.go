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

func (p *ProductionSvc) QueuedJobs(ctx context.Context) map[uuid.UUID][]*game.Job {
	return p.storage.ProductJobsByTown(ctx, TownFromContext(ctx))
}

func (p *ProductionSvc) QueuedBuildings(ctx context.Context) []*game.Job {
	return p.storage.QueuedBuildings(ctx, TownFromContext(ctx))
}

func (p *ProductionSvc) CreateJob(ctx context.Context, inputJob *game.InputJob) error {
	var job = new(game.Job)
	job.InputJob = *inputJob
	job.ID = uuid.New()
	job.TownID = TownFromContext(ctx)
	job.Queued = time.Now().UTC()
	job.Completed = time.Now().Add(job.Hours).UTC()

	// check if we can make this job active
	if job.Type == game.JobTypeProduct {
		queuedJobs := p.QueuedJobs(ctx)
		if _, ok := queuedJobs[job.ProductJob.BuildingID]; !ok {
			// no jobs found for this building, make this job active.
			job.Status = game.JobStatusActive
			job.Started = time.Now().UTC()
		}
	}
	if job.Type == game.JobTypeBuilding {
		if len(p.QueuedBuildings(ctx)) == 0 {
			job.Status = game.JobStatusActive
			job.Started = time.Now().UTC()
		}
	}

	return p.storage.CreateJob(ctx, TownFromContext(ctx), job)
}

func (p *ProductionSvc) ProductJobsCompleted(ctx context.Context) map[uuid.UUID][]*game.Job {
	return p.storage.ProductJobsCompleted(ctx)
}

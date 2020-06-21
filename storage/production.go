package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type ProductionRepository struct {
	jobsByTown map[uuid.UUID][]uuid.UUID
	jobs       map[uuid.UUID]*game.Job
}

func NewProductionRepository() *ProductionRepository {
	return &ProductionRepository{
		jobsByTown: map[uuid.UUID][]uuid.UUID{},
		jobs:       map[uuid.UUID]*game.Job{},
	}
}

func (p *ProductionRepository) ProductJobsByTown(ctx context.Context, townID uuid.UUID) map[uuid.UUID][]*game.Job {
	jbt, ok := p.jobsByTown[townID]
	if !ok {
		return nil
	}

	var jobs = make(map[uuid.UUID][]*game.Job)
	for _, jb := range jbt {
		if j, ok := p.jobs[jb]; ok && j.Type == game.JobTypeProduct {
			jobs[j.ProductJob.BuildingID] = append(jobs[j.ProductJob.BuildingID], j)
		}
	}

	return jobs
}

func (p *ProductionRepository) QueuedBuildings(ctx context.Context, townID uuid.UUID) []*game.Job {
	jbt, ok := p.jobsByTown[townID]
	if !ok {
		return nil
	}

	var jobs []*game.Job
	for _, jb := range jbt {
		if j, ok := p.jobs[jb]; ok && j.Type == game.JobTypeBuilding {
			jobs = append(jobs, j)
		}
	}

	return jobs
}

func (p *ProductionRepository) CreateJob(ctx context.Context, townID uuid.UUID, job *game.Job) error {
	p.jobsByTown[townID] = append(p.jobsByTown[townID], job.ID)
	p.jobs[job.ID] = job

	return nil
}

func (p *ProductionRepository) ProductJobsCompleted(ctx context.Context) map[uuid.UUID][]*game.Job {
	var jobs = make(map[uuid.UUID][]*game.Job)
	for _, j := range p.jobs {
		if j.Type == game.JobTypeProduct && j.IsCompleted() {
			jobs[j.TownID] = append(jobs[j.TownID], j)
		}
	}

	return jobs
}

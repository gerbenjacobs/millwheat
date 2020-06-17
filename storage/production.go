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

func (p *ProductionRepository) QueuedJobs(ctx context.Context, townID uuid.UUID) []*game.Job {
	jbt, ok := p.jobsByTown[townID]
	if !ok {
		return nil
	}

	var jobs []*game.Job
	for _, jb := range jbt {
		if j, ok := p.jobs[jb]; ok {
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

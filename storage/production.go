package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/game/data"
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
		if j, ok := p.jobs[jb]; ok && j.Type == game.JobTypeProduct && j.Status != game.JobStatusCompleted {
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
		if j, ok := p.jobs[jb]; ok && j.Type == game.JobTypeBuilding && j.Status != game.JobStatusCompleted {
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
func (p *ProductionRepository) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status game.JobStatus) error {
	job, ok := p.jobs[jobID]
	if !ok {
		return errors.New("job not found")
	}

	job.Status = status
	p.jobs[jobID] = job

	return nil
}

func (p *ProductionRepository) CancelJob(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) error {
	jobs, ok := p.jobsByTown[townID]
	if !ok {
		return errors.New("jobs not found for town")
	}

	var found bool
	var newJobs []uuid.UUID
	for _, job := range jobs {
		if job == jobID {
			found = true
		} else {
			newJobs = append(newJobs, job)
		}
	}

	if !found {
		return errors.New("no matching job found in this town")
	}

	job, ok := p.jobs[jobID]
	if !ok {
		return errors.New("jobs not found")
	}

	if job.Status == game.JobStatusCompleted {
		return errors.New("job is already completed")
	}

	delete(p.jobs, jobID)
	p.jobsByTown[townID] = newJobs

	return nil
}

func (p *ProductionRepository) RevertJobResources(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) ([]game.ItemSet, error) {
	job, ok := p.jobs[jobID]
	if !ok {
		return nil, errors.New("job not found")
	}

	var items []game.ItemSet
	switch job.Type {
	case game.JobTypeProduct:
		items = job.ProductJob.Consumption
	case game.JobTypeBuilding:
		building, err := game.CreateBuilding(data.Buildings[job.BuildingJob.Type], job.BuildingJob.Level)
		if err != nil {
			return nil, errors.New("failed to find building")
		}
		items = building.Consumption
	}

	return items, nil
}

func (p *ProductionRepository) JobsCompleted(ctx context.Context) map[uuid.UUID][]*game.Job {
	var jobs = make(map[uuid.UUID][]*game.Job)
	for _, j := range p.jobs {
		if j.ReadyForCompletion() {
			jobs[j.TownID] = append(jobs[j.TownID], j)
		}
	}

	return jobs
}

func (p *ProductionRepository) ReshuffleQueue(ctx context.Context, townID uuid.UUID) {
	for _, j := range p.oldestQueuedJobs(townID) {
		if j != nil {
			j.Status = game.JobStatusActive
			j.Started = time.Now().UTC()
			j.Completed = j.Started.Add(j.Duration)
			p.jobs[j.ID] = j
		}
	}
}

func (p *ProductionRepository) oldestQueuedJobs(townID uuid.UUID) []*game.Job {
	jobs, ok := p.jobsByTown[townID]
	if !ok {
		return nil
	}

	var hasBuildingInProduction = false
	var oldestBuildingJob *game.Job = nil

	var hasProductInProduction = make(map[uuid.UUID]bool)
	var oldestProductJobs = make(map[uuid.UUID]*game.Job)
	for _, jobID := range jobs {
		if job, ok := p.jobs[jobID]; ok {
			// Products
			if job.Type == game.JobTypeProduct && job.Status == game.JobStatusActive {
				hasProductInProduction[job.BuildingID()] = true
			}
			if _, ok := hasProductInProduction[job.BuildingID()]; ok {
				// Skip looking for products in this building, already in progress
				continue
			}

			if job.Type == game.JobTypeProduct && job.Status == game.JobStatusQueued {
				currentOldest, ok := oldestProductJobs[job.BuildingID()]
				if ok {
					if currentOldest.Queued.After(job.Queued) {
						oldestProductJobs[job.BuildingID()] = job
					}
				} else {
					oldestProductJobs[job.BuildingID()] = job
				}
			}

			// Buildings
			if job.Type == game.JobTypeBuilding && job.IsActive() {
				// Skip looking for buildings, already 1 active
				hasBuildingInProduction = true
				continue
			}
			if !hasBuildingInProduction && job.Type == game.JobTypeBuilding && job.Status == game.JobStatusQueued {
				if oldestBuildingJob == nil || oldestBuildingJob.Queued.After(job.Queued) {
					oldestBuildingJob = job
				}
			}
		}
	}

	var oldestJobs []*game.Job
	if oldestBuildingJob != nil {
		oldestJobs = append(oldestJobs, oldestBuildingJob)
	}
	for _, j := range oldestProductJobs {
		oldestJobs = append(oldestJobs, j)
	}

	return oldestJobs
}

package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/game/data"
)

var CacheDurationJobs = 6 * time.Hour

type ProductionRepository struct {
	db             *sql.DB
	jobCache       *cache.Cache
	jobByTownCache *cache.Cache
}

func NewProductionRepository(db *sql.DB) *ProductionRepository {
	jc := cache.New(CacheDurationJobs, time.Hour)
	jtc := cache.New(CacheDurationJobs, time.Hour)

	return &ProductionRepository{
		db:             db,
		jobCache:       jc,
		jobByTownCache: jtc,
	}
}

func (p *ProductionRepository) jobsByTown(ctx context.Context, townID uuid.UUID) (jbt []uuid.UUID, err error) {
	tc, ok := p.jobByTownCache.Get(townID.String())
	if !ok {
		jbt, err = p.getJobsByTownFromDatabase(ctx, townID)
		if err != nil {
			return nil, err
		}
	} else {
		jbt, _ = tc.([]uuid.UUID)
	}

	return jbt, err
}

func (p *ProductionRepository) jobByID(ctx context.Context, jobID uuid.UUID) (job *game.Job, err error) {
	tc, ok := p.jobCache.Get(jobID.String())
	if !ok {
		job, err = p.getJobFromDatabase(ctx, jobID)
		if err != nil {
			return nil, err
		}
	} else {
		job, _ = tc.(*game.Job)
	}

	return job, err
}

func (p *ProductionRepository) ProductJobsByTown(ctx context.Context, townID uuid.UUID) map[uuid.UUID][]*game.Job {
	jbt, err := p.jobsByTown(ctx, townID)
	if err != nil {
		logrus.Errorf("failed to get jobs by town: %s", err)
		return nil
	}

	var jobs = make(map[uuid.UUID][]*game.Job)
	for _, jb := range jbt {
		j, err := p.jobByID(ctx, jb)
		if err != nil {
			continue
		}

		if j.Type == game.JobTypeProduct && j.Status != game.JobStatusCompleted {
			jobs[j.ProductJob.BuildingID] = append(jobs[j.ProductJob.BuildingID], j)
		}
	}

	return jobs
}

func (p *ProductionRepository) QueuedBuildings(ctx context.Context, townID uuid.UUID) []*game.Job {
	jbt, err := p.jobsByTown(ctx, townID)
	if err != nil {
		logrus.Errorf("failed to get jobs by town: %s", err)
		return nil
	}

	var jobs []*game.Job
	for _, jb := range jbt {
		j, err := p.jobByID(ctx, jb)
		if err != nil {
			continue
		}
		if j.Type == game.JobTypeBuilding && j.Status != game.JobStatusCompleted {
			jobs = append(jobs, j)
		}
	}

	return jobs
}

func (p *ProductionRepository) CreateJob(ctx context.Context, townID uuid.UUID, job *game.Job) error {
	return p.addJobToDatabase(ctx, townID, job)
}
func (p *ProductionRepository) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status game.JobStatus) error {
	return p.updateJobStatusInDatabase(ctx, jobID, status)
}

func (p *ProductionRepository) CancelJob(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) error {
	jobs, err := p.jobsByTown(ctx, townID)
	if err != nil {
		logrus.Errorf("failed to get jobs by town: %s", err)
		return nil
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

	job, err := p.jobByID(ctx, jobID)
	if err != nil {
		return err
	}

	if job.Status == game.JobStatusCompleted {
		return errors.New("job is already completed")
	}

	return p.deleteJobFromDatabase(ctx, townID, jobID)
}

func (p *ProductionRepository) RevertJobResources(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) ([]game.ItemSet, error) {
	job, err := p.jobByID(ctx, jobID)
	if err != nil {
		return nil, err
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
	jobs, err := p.getCompletedJobsFromDatabase(ctx)
	if err != nil {
		logrus.Errorf("failed to get completed jobs: %s", err)
		return nil
	}

	return jobs
}

func (p *ProductionRepository) ReshuffleQueue(ctx context.Context, townID uuid.UUID) {
	for _, j := range p.oldestQueuedJobs(townID) {
		if j != nil {
			j.Status = game.JobStatusActive
			j.Started = time.Now().UTC()
			j.Completed = j.Started.Add(j.Duration)
			if err := p.updateJobInDatabase(ctx, j); err != nil {
				logrus.Errorf("failed to update job in database: %s", err)
			}
		}
	}
}

func (p *ProductionRepository) oldestQueuedJobs(townID uuid.UUID) []*game.Job {
	jobs, err := p.jobsByTown(context.Background(), townID)
	if err != nil {
		logrus.Errorf("failed to get jobs by town: %s", err)
		return nil
	}

	var hasBuildingInProduction = false
	var oldestBuildingJob *game.Job = nil

	var hasProductInProduction = make(map[uuid.UUID]bool)
	var oldestProductJobs = make(map[uuid.UUID]*game.Job)
	for _, jobID := range jobs {
		job, err := p.jobByID(context.Background(), jobID)
		if err != nil {
			continue
		}

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

	var oldestJobs []*game.Job
	if oldestBuildingJob != nil {
		oldestJobs = append(oldestJobs, oldestBuildingJob)
	}
	for _, j := range oldestProductJobs {
		oldestJobs = append(oldestJobs, j)
	}

	return oldestJobs
}

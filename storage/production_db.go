package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
)

func (p *ProductionRepository) getJobsByTownFromDatabase(ctx context.Context, townID uuid.UUID) ([]uuid.UUID, error) {
	tid, _ := townID.MarshalBinary()

	rows, err := p.db.QueryContext(ctx, "SELECT id FROM jobs WHERE townId = ? ORDER BY queued", tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, id)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	p.jobByTownCache.Set(townID.String(), jobs, CacheDurationJobs)

	return jobs, nil
}

func (p *ProductionRepository) getJobFromDatabase(ctx context.Context, jobID uuid.UUID) (*game.Job, error) {
	tid, _ := jobID.MarshalBinary()
	row := p.db.QueryRowContext(ctx, "SELECT id, townId, type, jobData, queued, started, completed, status FROM jobs WHERE id = ?", tid)

	var job game.Job
	var jobData []byte
	err := row.Scan(&job.ID, &job.TownID, &job.Type, &jobData, &job.Queued, &job.Started, &job.Completed, &job.Status)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("job with ID %q not found", jobID)
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning for jobs: %v", err)
	}

	if err := json.Unmarshal(jobData, &job.InputJob); err != nil {
		return nil, err
	}

	p.jobCache.Set(jobID.String(), &job, CacheDurationJobs)

	return &job, nil
}

func (p *ProductionRepository) addJobToDatabase(ctx context.Context, townID uuid.UUID, job *game.Job) error {
	jobID, _ := job.ID.MarshalBinary()
	tid, _ := townID.MarshalBinary()
	jobData, err := json.Marshal(job.InputJob)
	if err != nil {
		return err
	}

	// write building to database
	stmt, err := p.db.PrepareContext(ctx, "INSERT INTO jobs (id, townId, type, jobData, queued, started, completed, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, jobID, tid, job.Type, jobData, job.Queued, job.Started, job.Completed, job.Status)
	if err != nil {
		return err
	}

	// update caches
	jbt, err := p.jobsByTown(ctx, townID)
	if err != nil {
		logrus.Errorf("failed to fetch job by town cache")
	}
	jbt = append(jbt, job.ID)
	p.jobByTownCache.Set(townID.String(), jbt, CacheDurationJobs)
	p.jobCache.Set(job.ID.String(), job, CacheDurationJobs)

	return nil
}

func (p *ProductionRepository) updateJobStatusInDatabase(ctx context.Context, jobID uuid.UUID, status game.JobStatus) error {
	job, err := p.jobByID(ctx, jobID)
	if err != nil {
		return err
	}
	jid, _ := jobID.MarshalBinary()

	query := "UPDATE jobs SET status = ?  WHERE id = ?"
	_, err = p.db.ExecContext(ctx, query, status, jid)
	if err != nil {
		return err
	}

	// update job struct and save to cache
	job.Status = status
	p.jobCache.Set(job.ID.String(), job, CacheDurationJobs)

	return nil
}

func (p *ProductionRepository) updateJobInDatabase(ctx context.Context, job *game.Job) error {
	jid, _ := job.ID.MarshalBinary()
	jobData, err := json.Marshal(job.InputJob)
	if err != nil {
		return err
	}

	query := "UPDATE jobs SET jobData = ?, queued = ?, started = ?, completed = ?, status = ? WHERE id = ?"
	_, err = p.db.ExecContext(ctx, query, jobData, job.Queued, job.Started, job.Completed, job.Status, jid)
	if err != nil {
		return err
	}

	// update job struct and save to cache
	p.jobCache.Set(job.ID.String(), job, CacheDurationJobs)

	return nil
}

func (p *ProductionRepository) deleteJobFromDatabase(ctx context.Context, townID uuid.UUID, jobID uuid.UUID) error {
	jid, _ := jobID.MarshalBinary()

	query := "DELETE FROM jobs WHERE id = ?"
	_, err := p.db.ExecContext(ctx, query, jid)
	if err != nil {
		return err
	}

	// update job struct and save to cache
	jbt, err := p.jobsByTown(ctx, townID)
	if err != nil {
		logrus.Errorf("failed to fetch job by town cache")
	}
	var newJbt []uuid.UUID
	for _, i := range jbt {
		if i != jobID {
			newJbt = append(newJbt, i)
		}
	}
	p.jobCache.Delete(jobID.String())
	p.jobByTownCache.Set(townID.String(), newJbt, CacheDurationJobs)
	return nil
}

func (p *ProductionRepository) getCompletedJobsFromDatabase(ctx context.Context) (map[uuid.UUID][]*game.Job, error) {
	q := "SELECT id, townId, type, jobData, queued, started, completed, status FROM jobs WHERE status = ? AND completed <= ?"
	rows, err := p.db.QueryContext(ctx, q, game.JobStatusActive, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs = make(map[uuid.UUID][]*game.Job)
	for rows.Next() {
		var job game.Job
		var jobData []byte
		err := rows.Scan(&job.ID, &job.TownID, &job.Type, &jobData, &job.Queued, &job.Started, &job.Completed, &job.Status)
		switch {
		case err != nil:
			return nil, fmt.Errorf("unknown error while scanning: %v", err)
		}

		if err := json.Unmarshal(jobData, &job.InputJob); err != nil {
			return nil, err
		}

		jobs[job.TownID] = append(jobs[job.TownID], &job)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

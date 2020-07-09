package game

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	JobTypeProduct JobType = iota
	JobTypeBuilding
)
const (
	JobStatusQueued JobStatus = iota
	JobStatusActive
	JobStatusCompleted
)

type JobType int
type JobStatus int

type Job struct {
	ID     uuid.UUID
	TownID uuid.UUID
	InputJob
	Queued    time.Time
	Started   time.Time
	Completed time.Time
	Status    JobStatus
}

type InputJob struct {
	Type        JobType
	ProductJob  *ProductJob
	BuildingJob *BuildingJob
	Duration    time.Duration
}

type ProductJob struct {
	BuildingID  uuid.UUID
	Consumption ItemSetSlice
	Production  ItemSetSlice
}
type BuildingJob struct {
	ID    uuid.UUID
	Type  BuildingType
	Level int
}

func (j *Job) String() string {
	return fmt.Sprintf("[%s] (%d) %s  -- %s %s %s -- will take: %s", j.ID, j.Type, j.Status, j.QueuedAt(), j.StartedAt(), j.Completed.Format(time.RFC3339), j.Duration)
}

func (js JobStatus) String() string {
	switch js {
	case JobStatusQueued:
		return "Queued"
	case JobStatusActive:
		return "In progress"
	case JobStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

func (jt JobType) String() string {
	switch jt {
	case JobTypeProduct:
		return "Product"
	case JobTypeBuilding:
		return "Building"
	default:
		return "Unknown"
	}
}

func (j *Job) BuildingID() uuid.UUID {
	switch j.Type {
	case JobTypeBuilding:
		return j.BuildingJob.ID
	case JobTypeProduct:
		return j.ProductJob.BuildingID
	default:
		logrus.Warnf("unknown job type for job %s", j.ID)
		return uuid.UUID{}
	}
}

func (j *Job) QueuedAt() string {
	return j.Queued.Format("2006-01-02 15:04")
}

func (j *Job) StartedAt() string {
	return j.Started.Format("2006-01-02 15:04")
}

func (j *Job) ReadyAt() string {
	ready := j.Started.Add(j.Duration)
	return ready.Format("2006-01-02 15:04")
}

func (j *Job) Progress() int {
	if j.Status != JobStatusActive {
		return 0
	}

	completed := j.Completed.Sub(time.Now().UTC()).Minutes()
	if completed <= 0 {
		return 100
	}

	p := 100 - (completed/j.Duration.Minutes())*100
	return int(math.Floor(p))
}

func (j *Job) ReadyForCompletion() bool {
	return j.IsActive() && j.Completed.Before(time.Now().UTC())
}

func (j *Job) IsActive() bool {
	return j.Status == JobStatusActive
}

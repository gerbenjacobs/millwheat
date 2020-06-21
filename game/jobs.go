package game

import (
	"fmt"
	"math"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
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
	Hours       time.Duration
}

type ProductJob struct {
	BuildingID  uuid.UUID
	Consumption []ItemSet
	Production  []ItemSet
}
type BuildingJob struct {
	ID    uuid.UUID
	Type  BuildingType
	Level int
}

func (j *Job) String() string {
	return fmt.Sprintf("[%s] (%d) %s - created at: %s -- will take: %s", j.ID, j.Type, spew.Sdump(j.ProductJob), j.Queued, j.Hours)
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

func (j *Job) QueuedAt() string {
	return j.Queued.Format("2006-01-02 15:04")
}

func (j *Job) StartedAt() string {
	return j.Started.Format("2006-01-02 15:04")
}

func (j *Job) ReadyAt() string {
	ready := j.Started.Add(j.Hours)
	return ready.Format("2006-01-02 15:04")
}

func (j *Job) Progress() int {
	if j.Status != JobStatusActive {
		return 0
	}

	completed := j.Completed.Sub(time.Now().UTC()).Minutes()
	if completed == 0 {
		return 0
	}

	p := 100 - (completed/j.Hours.Minutes())*100
	return int(math.Floor(p))
}

func (j *Job) IsCompleted() bool {
	if j.Status != JobStatusActive {
		return false
	}

	return j.Completed.Before(time.Now().UTC())
}

func (j *Job) IsActive() bool {
	return j.Status == JobStatusActive
}

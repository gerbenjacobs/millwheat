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

type JobType int

type Job struct {
	ID          uuid.UUID
	Type        JobType
	ProductJob  *ProductJob
	BuildingJob *BuildingJob
	Created     time.Time
	Completed   time.Time
	Hours       time.Duration
	Active      bool
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
	return fmt.Sprintf("[%s] (%d) %s - created at: %s -- will take: %s", j.ID, j.Type, spew.Sdump(j.ProductJob), j.Created, j.Hours)
}

func (j *Job) CreatedAt() string {
	return j.Created.Format("2006-01-02 15:04")
}

func (j *Job) ReadyAt() string {
	ready := j.Created.Add(j.Hours)
	return ready.Format("2006-01-02 15:04")
}

func (j *Job) Progress() int {
	if !j.Active {
		return 0
	}

	completed := j.Completed.Sub(time.Now().UTC()).Minutes()
	if completed == 0 {
		return 0
	}

	p := 100 - (completed/j.Hours.Minutes())*100
	return int(math.Floor(p))
}

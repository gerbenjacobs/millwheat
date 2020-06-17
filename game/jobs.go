package game

import (
	"fmt"
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
	return fmt.Sprintf("[%s] (%d) %s - created at: %s", j.ID, j.Type, spew.Sdump(j.ProductJob), j.Created)
}

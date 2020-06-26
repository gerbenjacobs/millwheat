package handler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/services"
)

func (h *Handler) Tick(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)

	tickHandler := func() {
		h.evaluateJobs(ctx)
	}

	// initial tick run
	tickHandler()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _ = <-t.C:
				tickHandler()
			}
		}
	}()
}

func (h *Handler) evaluateJobs(ctx context.Context) {
	completedJobs := h.ProductionSvc.JobsCompleted(ctx)
	if len(completedJobs) > 0 {
		logrus.Debug("handling completed jobs")
	}

	for townID, jobs := range completedJobs {
		for _, job := range jobs {
			ctx = context.WithValue(ctx, services.CtxKeyTownID, townID)
			var err error
			switch job.Type {
			case game.JobTypeProduct:
				err = h.TownSvc.GiveToWarehouse(ctx, job.ProductJob.Production)
			case game.JobTypeBuilding:
				if job.BuildingJob.Level == 1 {
					err = h.TownSvc.AddBuilding(ctx, job.BuildingJob.Type)
				} else {
					err = h.TownSvc.UpgradeBuilding(ctx, job.BuildingJob.ID)
				}
			}
			if err != nil {
				logrus.Errorf("failed to resolve job production: %s", err)
				continue
			}
			if err := h.ProductionSvc.UpdateJobStatus(ctx, job.ID, game.JobStatusCompleted); err != nil {
				logrus.Errorf("failed to update job status for %s: %s", job.ID, err)
			}

			// reshuffle the queue
			h.ProductionSvc.ReshuffleQueue(ctx)
		}
	}
}

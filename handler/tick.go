package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func (h *Handler) Tick(ctx context.Context) {
	t := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-t.C:
				fmt.Println("Tick at", t)
				h.evaluateJobs(ctx)
			}
		}
	}()
}

func (h *Handler) evaluateJobs(ctx context.Context) {
	jobs := h.ProductionSvc.ProductJobsCompleted(ctx)
	spew.Dump(jobs)
}

package handler

import (
	"context"
	"fmt"
	"time"
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
			}
		}
	}()
}

func (h *Handler) evaluateJobs() {

}

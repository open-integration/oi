package tasks

import (
	"context"
	"time"

	"github.com/open-integration/core/pkg/task"
	"github.com/open-integration/core/pkg/utils"
)

type (
	// ticker task
	ticker struct {
		meta         task.Metadata
		tickInterval time.Duration
		totalTime    time.Duration
	}
)

// Run implements the runner interface to run the task
func (t *ticker) Run(ctx context.Context, options task.RunOptions) ([]byte, error) {
	t.meta.Time.StartedAt = utils.TimeNow()
	ticker := time.NewTicker(t.tickInterval)
	finish := time.NewTimer(t.totalTime)
	for {
		select {
		// Finish the task execution
		case _ = <-finish.C:
			t.meta.Time.FinishedAt = utils.TimeNow()
			return nil, nil
		case _ = <-ticker.C:
			options.EventReporter.Report("tick", nil)
		}
	}
}

func (t *ticker) Metadata() task.Metadata {
	return t.meta
}

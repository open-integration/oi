package tasks

import (
	"context"
	"time"

	"github.com/open-integration/core/pkg/task"
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
	ticker := time.NewTicker(t.tickInterval)
	finish := time.NewTimer(t.totalTime)
	for {
		select {
		// Finish the task execution
		case _ = <-finish.C:
			return nil, nil
		case _ = <-ticker.C:
			options.EventReporter.Report("tick", nil)
		}
	}
}

func (t *ticker) Name() string {
	return t.meta.Name
}

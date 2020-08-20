package tasks

import (
	"context"
	"time"

	"github.com/open-integration/oi/core/task"
	"github.com/open-integration/oi/pkg/logger"
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
	lgr := logger.New(&logger.Options{
		FilePath: options.FD.File(),
	})
	for {
		select {
		// Finish the task execution
		case _ = <-finish.C:
			return nil, nil
		case _ = <-ticker.C:
			if err := options.EventReporter.Report("tick", nil); err != nil {
				lgr.Error("Failed to report event", "event", "tick", "err", err.Error())
			}
		}
	}
}

func (t *ticker) Name() string {
	return t.meta.Name
}

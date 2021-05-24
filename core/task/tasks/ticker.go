package tasks

import (
	"context"
	"fmt"
	"io"
	"os"
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
	f, err := os.OpenFile(options.FD.File(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, fmt.Errorf("failed to open file descriptor %s: %w", options.FD.File(), err)
	}
	lgr := logger.New(logger.Options{
		WriterHandlers: []io.Writer{f},
	})
	for {
		select {
		// Finish the task execution
		case <-finish.C:
			return nil, nil
		case <-ticker.C:
			if err := options.EventReporter.Report("tick", nil); err != nil {
				lgr.Info("Failed to report event", "event", "tick", "err", err.Error())
			}
		}
	}
}

func (t *ticker) Name() string {
	return t.meta.Name
}

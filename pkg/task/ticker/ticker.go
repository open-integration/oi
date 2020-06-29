package ticker

import (
	"time"

	"github.com/open-integration/core/pkg/event/reporter"
	"github.com/open-integration/core/pkg/task"
)

type (
	// Ticker task
	Ticker struct {
		TickInterval time.Duration
		TotalTime    time.Duration
	}

	// Options to create new ticker task
	// Name by default is "ticker"
	Options struct {
		Name         string
		TickInterval time.Duration
		TotalTime    time.Duration
	}
)

// New creates new Ticker task
func New(opt Options) task.Task {
	name = "ticker"
	if opt.Name != "" {
		name = opt.Name
	}
	return task.Task{
		Metadata: task.Metadata{
			Name: name,
		},
		Runner: Ticker{
			opt.TickInterval,
			opt.TotalTime,
		},
	}
}

// Run implements the runner interface to run the task
func (t Ticker) Run(eventReporter reporter.EventReporter) ([]byte, error) {
	ticker := time.NewTicker(t.TickInterval)
	finish := time.NewTimer(t.TotalTime)
	for {
		select {
		// Finish the task execution
		case _ = <-finish.C:
			return nil, nil
		case _ = <-ticker.C:
			eventReporter.Report("tick", nil)
		}
	}
	return nil, nil
}

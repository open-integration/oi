package tasks

import (
	"context"

	"github.com/open-integration/core/pkg/task"
	"github.com/open-integration/core/pkg/utils"
)

type (
	// function task will run in same process
	function struct {
		meta     task.Metadata
		function func(context.Context, task.RunOptions) ([]byte, error)
	}
)

// Run implements the runner interface to run the task
func (f *function) Run(ctx context.Context, options task.RunOptions) ([]byte, error) {
	f.meta.Time.StartedAt = utils.TimeNow()
	return f.function(ctx, options)
}

func (f *function) Metadata() task.Metadata {
	return f.meta
}

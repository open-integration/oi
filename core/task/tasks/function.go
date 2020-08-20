package tasks

import (
	"context"

	"github.com/open-integration/core/core/task"
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
	res, err := f.function(ctx, options)
	return res, err
}

func (f *function) Name() string {
	return f.meta.Name
}

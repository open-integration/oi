package tasks

import (
	"context"
	"time"

	"github.com/open-integration/open-integration/core/task"
)

// NewSerivceTask build task task calls a service with arguments
func NewSerivceTask(name string, service string, endpoint string, args ...task.Argument) task.Task {
	return &svc{
		meta: task.Metadata{
			Name: name,
		},
		endpoint:  endpoint,
		name:      service,
		arguments: args,
	}
}

// NewTickerTask builds task that will send event every tickInterval till it stops on totalTime
func NewTickerTask(name string, tickInterval time.Duration, totalTime time.Duration) task.Task {
	return &ticker{
		meta: task.Metadata{
			Name: name,
		},
		tickInterval: tickInterval,
		totalTime:    totalTime,
	}
}

// NewFunctionTask build task that will be executed in same process
func NewFunctionTask(name string, fn func(context.Context, task.RunOptions) ([]byte, error)) task.Task {
	return &function{
		meta: task.Metadata{
			Name: name,
		},
		function: fn,
	}
}

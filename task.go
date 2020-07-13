package core

import (
	"time"

	"github.com/open-integration/core/pkg/task"
	"github.com/open-integration/core/pkg/task/tasks"
)

// NewSerivceTask build task task calls a service with arguments
func NewSerivceTask(name string, service string, endpoint string, arg ...task.Argument) task.Task {
	return tasks.NewSerivceTask(name, service, endpoint, arg...)
}

// NewTickerTask builds task that will send event every tickInterval till it stops on totalTime
func NewTickerTask(name string, tickInterval time.Duration, totalTime time.Duration) task.Task {
	return tasks.NewTickerTask(name, tickInterval, totalTime)
}

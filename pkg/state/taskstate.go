package state

import (
	"time"

	"github.com/open-integration/core/pkg/task"
)

type (
	// TaskState is a representation of a state of one task
	TaskState struct {
		State  string    `yaml:"state"`
		Status string    `yaml:"status"`
		Task   task.Task `yaml:"task"`
		Times  TaskTimes `yaml:"times"`
		Output []byte    `yaml:"output"`
		Error  error     `yaml:"error"`
		Logger string    `yaml:"logger"`
	}

	TaskTimes struct {
		Started  time.Time `yaml:"started"`
		Finished time.Time `yaml:"finished"`
	}
)

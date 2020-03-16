package state

import (
	"github.com/open-integration/core/pkg/task"
)

type (
	// TaskState is a representation of a state of one task
	TaskState struct {
		State  string    `yaml:"state"`
		Status string    `yaml:"status"`
		Task   task.Task `yaml:"task"`
		Output string    `yaml:"output"`
		Error  error     `yaml:"error"`
		Logger string    `yaml:"logger"`
	}
)

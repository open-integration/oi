package task

import "time"

type (
	// Task is a task a pipeline should execute
	Task struct {
		Metadata Metadata
		Spec     Spec
	}

	// Metadata holds all the metadata of a pipeline
	Metadata struct {
		Name string
		Time struct {
			// StaredAt the time when the step was started execution
			StartedAt time.Time
			// FinishedAt the time when the step finished execution
			FinishedAt time.Time
		}
	}

	// Spec is the spec of a task
	Spec struct {
		Service   string
		Endpoint  string
		Arguments []Argument
		Cmd       *Command
	}

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value interface{}
	}

	// Command is a command to be executed as a sub-process
	Command struct {
		EnvironmentVariables []Argument
		// Program is executable program
		Program   string
		Arguments []string
		// WorkingDirectory default is the current directory
		WorkingDirectory string
	}
)

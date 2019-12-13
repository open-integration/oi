package core

type (
	// Task is a task a pipeline should execute
	Task struct {
		Metadata TaskMetadata
		Spec     TaskSpec
		SpecFunc func(*State) (*TaskSpec, error)
		// Condition a set of conditions to run before task execution
		Condition *Condition
	}

	// TaskMetadata holds all the metadata of a pipeline
	TaskMetadata struct {
		Name string
		// Reusable set to true to run the task multiple times
		Reusable bool
	}

	// TaskSpec is the spec of a task
	TaskSpec struct {
		Command              string
		EnvironmentVariables []string
		Detached             bool
		WorkingDirectory     string
		Path                 string
		Service              string
		Endpoint             string
		Arguments            []Argument
	}
)

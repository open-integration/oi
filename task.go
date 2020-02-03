package core

type (
	// Task is a task a pipeline should execute
	Task struct {
		Metadata TaskMetadata
		Spec     TaskSpec
	}

	// TaskMetadata holds all the metadata of a pipeline
	TaskMetadata struct {
		Name string
		// Reusable set to true to run the task multiple times
		Reusable bool
	}

	// TaskSpec is the spec of a task
	TaskSpec struct {
		Service   string
		Endpoint  string
		Arguments []Argument
	}
)

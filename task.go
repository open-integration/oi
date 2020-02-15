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
	}

	// TaskSpec is the spec of a task
	TaskSpec struct {
		Service   string
		Endpoint  string
		Arguments []Argument
	}
)

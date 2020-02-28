package task

type (
	// Task is a task a pipeline should execute
	Task struct {
		Metadata Metadata
		Spec     Spec
	}

	// Metadata holds all the metadata of a pipeline
	Metadata struct {
		Name string
	}

	// Spec is the spec of a task
	Spec struct {
		Service   string
		Endpoint  string
		Arguments []Argument
	}

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value interface{}
	}
)

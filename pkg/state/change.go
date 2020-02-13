package state

type (
	// Change is a sub state, represent the change in the state caused
	Change struct {
		Cause          string         `yaml:"cause"`
		TaskChanges    TaskChanges    `yaml:"task-changes"`
		PipelineChange PipelineChange `yaml:"pipeline-changes"`
	}

	// TaskChanges holds all the change a task get have
	TaskChanges struct {
		ID        string `yaml:"id"`
		Name      string `yaml:"name"`
		State     string `yaml:"state"`
		Status    string `yaml:"status"`
		EventID   string `yaml:"event-id"`
		LoggerID  string `yaml:"logger-id"`
		Arguments string `yaml:"arguments"`
		Output    string `yaml:"output"`
		Error     string `yaml:"error"`
	}
	// PipelineChange holds all the change a pipeline get have
	PipelineChange struct {
		States string `yaml:"state"`
	}
)

package state

type (
	// TaskState is a representation of a state of one task
	TaskState struct {
		ID        string `yaml:"id"`
		State     string `yaml:"state"`
		Status    string `yaml:"status"`
		Task      string `yaml:"task"`
		Output    string `yaml:"output"`
		Arguments string `yaml:"arguments"`
		Error     string `yaml:"error"`
		EventID   string `yaml:"event-id"`
		Logger    string `yaml:"logger"`
	}
)

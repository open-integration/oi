package event

import "time"

type (

	// Event - means that something happen
	Event struct {
		Metadata     Metadata               `yaml:"metadata"`
		RelatedTasks []string               `yaml:"related-tasks"`
		Payload      map[string]interface{} `yaml:"payload"`
	}

	// Metadata of a task
	Metadata struct {
		Name      string    `yaml:"name"`
		ID        string    `yaml:"id"`
		CreatedAt time.Time `yaml:"created-at"`
		Task      string    `yaml:"task"`
	}
)

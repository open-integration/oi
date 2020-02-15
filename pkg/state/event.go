package state

import "time"

const (
	EventEngineStarted  = "engine.started"
	EventEngineFinished = "engine.finished"
	EventTaskStarted    = "task.started"
	EventTaskFinished   = "task.finished"
)

type (

	// Event - means that something happen
	Event struct {
		Metadata     EventMetadata          `yaml:"metadata"`
		RelatedTasks []string               `yaml:"related-tasks"`
		Payload      map[string]interface{} `yaml:"payload"`
	}

	EventMetadata struct {
		Name      string
		ID        string
		CreatedAt time.Time
	}
)

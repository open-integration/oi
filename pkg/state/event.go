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
		Metadata EventMetadata
		Payload  map[string]interface{}
	}

	EventMetadata struct {
		Name      string
		ID        string
		CreatedAt time.Time
	}
)

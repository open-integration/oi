package openc

import "time"

const (
	EventEngineStarted  = "engine.started"
	EventEngineFinished = "engine.finished"
	EventTaskStarted    = "task.started"
	EventTaskFinished   = "task.finished"
)

type (
	event int

	// Event - means that something happen
	Event struct {
		Metadata EventMetadata
		Payload  map[string]interface{}
	}

	EventMetadata struct {
		Name      string
		ID        ID
		CreatedAt time.Time
	}
)

package state

import (
	"time"

	"github.com/open-integration/oi/core/task"
)

type (
	// UpdateRequest request update the state
	UpdateRequest struct {
		Metadata                     UpdateRequestMetadata
		ElectTasksRequest            *ElectTasksRequest
		AddRealtedTaskToEventReuqest *AddRealtedTaskToEventReuqest
		UpdateTaskStateRequest       *UpdateTaskStateRequest
		UpdateStateMetadataRequest   *UpdateStateMetadataRequest
	}

	// UpdateRequestMetadata metadata to UpdateRequest
	UpdateRequestMetadata struct {
		CreatedAt time.Time
	}

	// ElectTasksRequest request to mark tasks as elected
	ElectTasksRequest struct {
		Tasks []task.Task
	}

	// UpdateTaskStateRequest request to update task state
	UpdateTaskStateRequest struct {
		State TaskState
	}

	// UpdateStateMetadataRequest request to update task state metadata
	UpdateStateMetadataRequest struct {
		State  string
		Status string
	}

	// AddRealtedTaskToEventReuqest request to add related task
	// to the event
	AddRealtedTaskToEventReuqest struct {
		EventID string
		Task    []string
	}
)

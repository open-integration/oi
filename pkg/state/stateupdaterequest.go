package state

import (
	"time"

	"github.com/open-integration/core/pkg/task"
)

type (
	StateUpdateRequest struct {
		Metadata                     StateUpdateRequestMetadata
		ElectTasksRequest            *ElectTasksRequest
		AddRealtedTaskToEventReuqest *AddRealtedTaskToEventReuqest
		UpdateTaskStateRequest       *UpdateTaskStateRequest
		UpdateStateMetadataRequest   *UpdateStateMetadataRequest
	}

	StateUpdateRequestMetadata struct {
		CreatedAt time.Time
	}

	ElectTasksRequest struct {
		Tasks []task.Task
	}

	UpdateTaskStateRequest struct {
		State TaskState
	}

	UpdateStateMetadataRequest struct {
		State  string
		Status string
	}

	AddRealtedTaskToEventReuqest struct {
		EventID string
		Task    []string
	}
)

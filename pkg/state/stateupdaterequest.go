package state

import "time"

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
		Tasks []string
	}

	UpdateTaskStateRequest struct {
		State TaskState
	}

	UpdateStateMetadataRequest struct {
		State string
	}

	AddRealtedTaskToEventReuqest struct {
		EventID string
		Task    string
	}
)

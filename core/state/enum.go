package state

const (

	// EngineStateInProgress pipeline in execution progress
	EngineStateInProgress string = "in-progress"
	// EngineStateFinished pipeline is finished execution
	EngineStateFinished string = "finished"

	// EngineStatusSuccess marks the engine as finished successfully
	EngineStatusSuccess = "Success"

	// EngineStatusFailed marks the engine as finished with error
	EngineStatusFailed = "failed"

	// TaskStateElected task is in progress
	TaskStateElected string = "elected"

	// TaskStateInProgress task is in progress
	TaskStateInProgress string = EngineStateInProgress

	// TaskStateFinished task is finished
	TaskStateFinished string = EngineStateFinished

	// TaskStatusSuccess set on the task status in case task was finished successfully
	TaskStatusSuccess = "Success"

	// TaskStatusFailed set on the task status in case task was finished with error
	TaskStatusFailed = "failed"
)

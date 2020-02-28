package state

const (

	// EngineStateInProgress pipeline in execution progress
	EngineStateInProgress string = "in-progress"
	// EngineStateFinished pipeline is finished execution
	EngineStateFinished string = "finished"

	// TaskStateElected task is in progress
	TaskStateElected string = "elected"

	// TaskStateInProgress task is in progress
	TaskStateInProgress string = EngineStateInProgress

	// TaskStateFinished task is finished
	TaskStateFinished string = EngineStateFinished

	// TaskStatusSuccess set one the task status in case task was finished successfully
	TaskStatusSuccess = "Success"
	// TaskStatusFailed set one the task status in case task was finished and failed
	TaskStatusFailed = "failed"
)

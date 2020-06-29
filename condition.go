package core

import (
	"github.com/open-integration/core/pkg/condition"
)

// ConditionEngineStarted returns the condition that is evaluated to true on engine.started event
func ConditionEngineStarted() condition.Condition {
	return condition.EngineStarted()
}

// ConditionTaskFinished returns the condition that is evaluated to true on task.finished event
// and the task is marked as finished in the state
func ConditionTaskFinished(task string) condition.Condition {
	return condition.TaskFinished(task)
}

// ConditionTaskFinishedWithStatus returns the condition that is evaluated to true on task.finished event
// and the task is marked as finished in the state
// and the status is as given
func ConditionTaskFinishedWithStatus(task string, status string) condition.Condition {
	return condition.TaskFinishedWithStatus(task, status)
}

// ConditionCombined returns the condition that is evaluated to true when all the conditions are true
func ConditionCombined(conditions ...condition.Condition) condition.Condition {
	return condition.Combined(conditions...)
}

// ConditionTaskEventReported return the condition that satisfied when task reported event
// in format {TASK_NAME}.{EVENT}
func ConditionTaskEventReported(name string) condition.Condition {
	return condition.TaskEventReported(name)
}

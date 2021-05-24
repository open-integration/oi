package oi

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/open-integration/oi/core/condition"
	"github.com/open-integration/oi/core/engine"
	"github.com/open-integration/oi/core/task"
	"github.com/open-integration/oi/core/task/tasks"
	"github.com/open-integration/oi/pkg/logger"
)

var exit = func(code int) { os.Exit(code) }

type (
	// EngineOptions to create new engine
	EngineOptions struct {
		Pipeline engine.Pipeline
		// LogsDirectory path where to store logs
		LogsDirectory string
		Kubeconfig    *engine.KubernetesOptions
		Logger        logger.Logger
	}
)

// NewEngine create new engine
func NewEngine(opt *EngineOptions) engine.Engine {
	e, err := engine.New(&engine.Options{
		Pipeline:      opt.Pipeline,
		Kubeconfig:    opt.Kubeconfig,
		Logger:        opt.Logger,
		LogsDirectory: opt.LogsDirectory,
	})
	dieOnError(err)
	return e
}

func dieOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		exit(1)
	}
}

// HandleEngineError prints the error in case the engine.Run was failed and exit
func HandleEngineError(err error) {
	if err != nil {
		fmt.Printf("Failed to execute pipeline\nError: %v", err)
		exit(1)
	}
}

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

// NewSerivceTask build task task calls a service with arguments
func NewSerivceTask(name string, service string, endpoint string, arg ...task.Argument) task.Task {
	return tasks.NewSerivceTask(name, service, endpoint, arg...)
}

// NewTickerTask builds task that will send event every tickInterval till it stops on totalTime
func NewTickerTask(name string, tickInterval time.Duration, totalTime time.Duration) task.Task {
	return tasks.NewTickerTask(name, tickInterval, totalTime)
}

// NewFunctionTask build task that will be executed in same process
func NewFunctionTask(name string, fn func(context.Context, task.RunOptions) ([]byte, error)) task.Task {
	return tasks.NewFunctionTask(name, fn)
}

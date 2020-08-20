package engine

import (
	"github.com/open-integration/core/core/condition"
	"github.com/open-integration/core/core/event"
	"github.com/open-integration/core/core/state"
	"github.com/open-integration/core/core/task"
)

type (
	// Pipeline is the pipeline representation
	Pipeline struct {
		Metadata PipelineMetadata
		Spec     PipelineSpec
	}

	// PipelineMetadata holds all the metadata of a pipeline
	PipelineMetadata struct {
		Name         string
		OS           string
		Architecture string
	}

	// PipelineSpec is the spec of a pipeline
	PipelineSpec struct {
		Reactions []EventReaction
		Services  []Service
	}

	// EventReaction is a binding of an event to a function that builds tasks
	EventReaction struct {
		Condition condition.Condition
		Reaction  func(ev event.Event, state state.State) []task.Task
	}
)

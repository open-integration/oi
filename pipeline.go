package core

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
		Condition func(ev Event, state State) bool
		Reaction  func(ev Event, state State) []Task
	}
)

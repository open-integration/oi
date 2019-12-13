package openc

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
		Tasks    []Task
		Services []Service
	}
)

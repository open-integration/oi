package command

import (
	"encoding/json"
	"io/ioutil"

	"github.com/open-integration/oi/pkg/logger"
)

// UnmarshalPipeline unmarshals json into pipeline
func UnmarshalPipeline(data []byte) (r *Pipeline, err error) {
	err = json.Unmarshal(data, &r)
	return r, err
}

// Marshal marshals pipeline into json
func (r *Pipeline) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Metadata of pipeline
type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

// Pipeline structure to generate pipeline code
type Pipeline struct {
	Metadata       Metadata          `json:"metadata" yaml:"metadata"`
	Services       []PipelineService `json:"services" yaml:"services"`
	EventReactions []EventReaction   `json:"eventReactions" yaml:"eventReactions"`
}

// PipelineService represents dependency service
type PipelineService struct {
	Name    string `json:"name" yaml:"name"`
	As      string `json:"as" yaml:"as"`
	Version string `json:"version" yaml:"version"`
}

// EventReaction maps condition to reactions
type EventReaction struct {
	Condition string     `json:"condition" yaml:"condition"`
	Reaction  []Reaction `json:"reactions" yaml:"reactions"`
}

// Reaction on oi.engine event
type Reaction struct {
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
}

// LoadFromPath loads pipeline from filesystem
func LoadFromPath(location string, log logger.Logger) (*Pipeline, error) {
	log.Debug("Loading pipeline", "path", location)
	f, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	log.Debug("File loaded")

	return UnmarshalPipeline(f)
}

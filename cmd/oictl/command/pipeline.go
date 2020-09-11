package command

import (
	"encoding/json"
	"io/ioutil"

	"github.com/open-integration/oi/pkg/logger"
)

func UnmarshalPipeline(data []byte) (r *Pipeline, err error) {
	err = json.Unmarshal(data, &r)
	return r, err
}

func (r *Pipeline) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type Pipeline struct {
	Metadata       Metadata          `json:"metadata" yaml:"metadata"`
	Services       []PipelineService `json:"services" yaml:"services"`
	EventReactions []EventReaction   `json:"eventReactions" yaml:"eventReactions"`
}

type PipelineService struct {
	Name    string `json:"name" yaml:"name"`
	As      string `json:"as" yaml:"as"`
	Version string `json:"version" yaml:"version"`
}

type EventReaction struct {
	Condition string     `json:"condition" yaml:"condition"`
	Reaction  []Reaction `json:"reactions" yaml:"reactions"`
}

type Reaction struct {
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
}

func LoadFromPath(location string, log logger.Logger) (*Pipeline, error) {
	log.Debug("Loading pipeline", "path", location)
	f, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	log.Debug("File loaded")

	return UnmarshalPipeline(f)
}

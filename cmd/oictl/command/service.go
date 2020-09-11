package command

import "gopkg.in/hairyhenderson/yaml.v2"

func UnmarshalService(data []byte) (Service, error) {
	var r Service
	err := yaml.Unmarshal(data, &r)
	return r, err
}

func (r *Service) Marshal() ([]byte, error) {
	return yaml.Marshal(r)
}

type Service struct {
	Name      string              `yaml:"name" json:"name"`
	Project   string              `yaml:"project" json:"project"`
	Version   string              `yaml:"version" json:"version"`
	Endpoints map[string]Endpoint `yaml:"endpoints" json:"endpoints"`
	Types     []string            `yaml:"types" json:"types"`
}

type Endpoint struct {
	Name      string `yaml:"name" json:"name"`
	Arguments string `yaml:"arguments" json:"arguments"`
	Returns   string `yaml:"returns" json:"returns"`
}

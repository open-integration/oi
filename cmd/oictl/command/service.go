package command

import "gopkg.in/hairyhenderson/yaml.v2"

// UnmarshalService unmarshals json into Service
func UnmarshalService(data []byte) (Service, error) {
	var r Service
	err := yaml.Unmarshal(data, &r)
	return r, err
}

// Marshal service into json
func (r *Service) Marshal() ([]byte, error) {
	return yaml.Marshal(r)
}

// Service structure to generate service code
type Service struct {
	Name      string              `yaml:"name" json:"name"`
	Project   string              `yaml:"project" json:"project"`
	Version   string              `yaml:"version" json:"version"`
	Endpoints map[string]Endpoint `yaml:"endpoints" json:"endpoints"`
	Types     []string            `yaml:"types" json:"types"`
}

// Endpoint of a service
type Endpoint struct {
	Name      string `yaml:"name" json:"name"`
	Arguments string `yaml:"arguments" json:"arguments"`
	Returns   string `yaml:"returns" json:"returns"`
}

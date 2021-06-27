// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    service, err := UnmarshalService(bytes)
//    bytes, err = service.Marshal()

package types

import "encoding/json"

func UnmarshalService(data []byte) (Service, error) {
	var r Service
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Service) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Service struct {
	Endpoints *Endpoints `json:"endpoints,omitempty"`
}

type Endpoints struct {
	Call *Call `json:"call,omitempty"`
}

type Call struct {
	Arguments *Arguments `json:"arguments,omitempty"`
	Returns   *Returns   `json:"returns,omitempty"`
}

type Arguments struct {
	Content *string           `json:"Content,omitempty"`
	Headers []ArgumentsHeader `json:"Headers,omitempty"`
	URL     *string           `json:"URL,omitempty"`
	Verb    *string           `json:"Verb,omitempty"`
}

type ArgumentsHeader struct {
	Name  *string `json:"Name,omitempty"`
	Value *string `json:"Value,omitempty"`
}

type Returns struct {
	Body    string          `json:"Body"`
	Headers []ReturnsHeader `json:"Headers"`
	Status  float64         `json:"Status"`
}

type ReturnsHeader struct {
	Name  *string `json:"Name,omitempty"`
	Value *string `json:"Value,omitempty"`
}

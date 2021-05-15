// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    callArguments, err := UnmarshalCallArguments(bytes)
//    bytes, err = callArguments.Marshal()

package call

import "encoding/json"

func UnmarshalCallArguments(data []byte) (CallArguments, error) {
	var r CallArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CallArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CallArguments struct {
	Content *string  `json:"Content,omitempty"`
	Headers []Header `json:"Headers"`
	URL     string   `json:"URL"`
	Verb    string   `json:"Verb"`
}

type Header struct {
	Name  *string `json:"Name,omitempty"`
	Value *string `json:"Value,omitempty"`
}

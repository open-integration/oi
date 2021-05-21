// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    callReturns, err := UnmarshalCallReturns(bytes)
//    bytes, err = callReturns.Marshal()

package call

import "encoding/json"

func UnmarshalCallReturns(data []byte) (CallReturns, error) {
	var r CallReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CallReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CallReturns struct {
	Body    string   `json:"Body"`
	Headers []Header `json:"Headers"`
	Status  int      `json:"Status"`
}

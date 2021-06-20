// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addcardReturns, err := UnmarshalAddcardReturns(bytes)
//    bytes, err = addcardReturns.Marshal()

package addcard

import "encoding/json"

type AddcardReturns map[string]interface{}

func UnmarshalAddcardReturns(data []byte) (AddcardReturns, error) {
	var r AddcardReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddcardReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

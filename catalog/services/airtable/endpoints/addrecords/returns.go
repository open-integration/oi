// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addRecordsReturns, err := UnmarshalAddRecordsReturns(bytes)
//    bytes, err = addRecordsReturns.Marshal()

package addrecords

import "encoding/json"

type AddRecordsReturns map[string]interface{}

func UnmarshalAddRecordsReturns(data []byte) (AddRecordsReturns, error) {
	var r AddRecordsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddRecordsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

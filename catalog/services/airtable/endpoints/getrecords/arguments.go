// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getRecordsArguments, err := UnmarshalGetRecordsArguments(bytes)
//    bytes, err = getRecordsArguments.Marshal()

package getrecords

import "encoding/json"

func UnmarshalGetRecordsArguments(data []byte) (GetRecordsArguments, error) {
	var r GetRecordsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetRecordsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetRecordsArguments struct {
	APIKey     string  `json:"APIKey"`            // Airtable API key
	DatabaseID string  `json:"DatabaseID"`        // Database ID
	Formula    *string `json:"Formula,omitempty"` // Formula
	TableName  string  `json:"TableName"`         // Table name
}

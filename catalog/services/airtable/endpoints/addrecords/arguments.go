// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addRecordsArguments, err := UnmarshalAddRecordsArguments(bytes)
//    bytes, err = addRecordsArguments.Marshal()

package addrecords

import "encoding/json"

func UnmarshalAddRecordsArguments(data []byte) (AddRecordsArguments, error) {
	var r AddRecordsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddRecordsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AddRecordsArguments struct {
	APIKey     string   `json:"APIKey"`     // Airtable API key
	DatabaseID string   `json:"DatabaseID"` // Database ID
	Records    []Record `json:"Records,omitempty"`
	TableName  string   `json:"TableName"` // Table name
}

type Record struct {
	CreatedTime *string                `json:"createdTime,omitempty"`
	Deleted     *bool                  `json:"deleted,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
	ID          *string                `json:"id,omitempty"`
	Typecast    *bool                  `json:"typecast,omitempty"`
}

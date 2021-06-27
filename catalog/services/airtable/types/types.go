// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addRecordsArguments, err := UnmarshalAddRecordsArguments(bytes)
//    bytes, err = addRecordsArguments.Marshal()
//
//    addRecordsReturns, err := UnmarshalAddRecordsReturns(bytes)
//    bytes, err = addRecordsReturns.Marshal()
//
//    getRecordsArguments, err := UnmarshalGetRecordsArguments(bytes)
//    bytes, err = getRecordsArguments.Marshal()
//
//    getRecordsReturns, err := UnmarshalGetRecordsReturns(bytes)
//    bytes, err = getRecordsReturns.Marshal()

package types

import "encoding/json"

func UnmarshalAddRecordsArguments(data []byte) (AddRecordsArguments, error) {
	var r AddRecordsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddRecordsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAddRecordsReturns(data []byte) (AddRecordsReturns, error) {
	var r AddRecordsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddRecordsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalGetRecordsArguments(data []byte) (GetRecordsArguments, error) {
	var r GetRecordsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetRecordsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalGetRecordsReturns(data []byte) (GetRecordsReturns, error) {
	var r GetRecordsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetRecordsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AddRecordsArguments struct {
	Auth    Auth     `json:"Auth"`
	Records []Record `json:"Records,omitempty"`
}

type Auth struct {
	APIKey     string `json:"APIKey"`     // Airtable API key
	DatabaseID string `json:"DatabaseID"` // Database ID
	TableName  string `json:"TableName"`  // Table name
}

type Record struct {
	CreatedTime *string                `json:"createdTime,omitempty"`
	Deleted     *bool                  `json:"deleted,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
	ID          *string                `json:"id,omitempty"`
	Typecast    *bool                  `json:"typecast,omitempty"`
}

type AddRecordsReturns struct {
	Records []Record `json:"Records,omitempty"`
}

type GetRecordsArguments struct {
	Auth    Auth    `json:"Auth"`
	Formula *string `json:"Formula,omitempty"` // Formula
}

type GetRecordsReturns struct {
	Records []Record `json:"Records,omitempty"`
}

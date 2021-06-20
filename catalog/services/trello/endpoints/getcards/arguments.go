// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getcardsArguments, err := UnmarshalGetcardsArguments(bytes)
//    bytes, err = getcardsArguments.Marshal()

package getcards

import "encoding/json"

func UnmarshalGetcardsArguments(data []byte) (GetcardsArguments, error) {
	var r GetcardsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetcardsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetcardsArguments struct {
	App   string `json:"App"`   // Trello Application ID
	Board string `json:"Board"` // Trello Board ID
	Token string `json:"Token"` // Trello Token
}

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addcardArguments, err := UnmarshalAddcardArguments(bytes)
//    bytes, err = addcardArguments.Marshal()

package addcard

import "encoding/json"

func UnmarshalAddcardArguments(data []byte) (AddcardArguments, error) {
	var r AddcardArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddcardArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AddcardArguments struct {
	App         string   `json:"App"`                   // Trello Application ID
	Board       string   `json:"Board"`                 // Trello board ID
	Description *string  `json:"Description,omitempty"` // Trello description to set on card
	Labels      []string `json:"Labels"`                // Trello labels to apply on card
	List        string   `json:"List"`                  // Trello list ID
	Name        string   `json:"Name"`                  // Trello card name
	Token       string   `json:"Token"`                 // Trello token
}

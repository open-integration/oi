// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    archivecardsArguments, err := UnmarshalArchivecardsArguments(bytes)
//    bytes, err = archivecardsArguments.Marshal()

package archivecards

import "encoding/json"

func UnmarshalArchivecardsArguments(data []byte) (ArchivecardsArguments, error) {
	var r ArchivecardsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ArchivecardsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ArchivecardsArguments struct {
	App     string   `json:"App"`     // Trello Application ID
	CardIDs []string `json:"CardIDs"` // IDs to archive
	Token   string   `json:"Token"`   // Trello Token
}

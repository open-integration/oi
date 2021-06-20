// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    archivecardsReturns, err := UnmarshalArchivecardsReturns(bytes)
//    bytes, err = archivecardsReturns.Marshal()

package archivecards

import "encoding/json"

type ArchivecardsReturns map[string]interface{}

func UnmarshalArchivecardsReturns(data []byte) (ArchivecardsReturns, error) {
	var r ArchivecardsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ArchivecardsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

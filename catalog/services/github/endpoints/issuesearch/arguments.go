// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    issueSearchArguments, err := UnmarshalIssueSearchArguments(bytes)
//    bytes, err = issueSearchArguments.Marshal()

package issuesearch

import "encoding/json"

func UnmarshalIssueSearchArguments(data []byte) (IssueSearchArguments, error) {
	var r IssueSearchArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IssueSearchArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IssueSearchArguments struct {
	Query string `json:"query"`
	Token string `json:"token"`
}

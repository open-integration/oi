// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    geIssueCommentArguments, err := UnmarshalGeIssueCommentArguments(bytes)
//    bytes, err = geIssueCommentArguments.Marshal()

package getissuecomments

import "encoding/json"

func UnmarshalGeIssueCommentArguments(data []byte) (GeIssueCommentArguments, error) {
	var r GeIssueCommentArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GeIssueCommentArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GeIssueCommentArguments struct {
	Issue float64 `json:"issue"`
	Owner string  `json:"owner"`
	Repo  string  `json:"repo"`
	Token string  `json:"token"`
}

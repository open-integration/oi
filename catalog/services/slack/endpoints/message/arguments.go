// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    messageArguments, err := UnmarshalMessageArguments(bytes)
//    bytes, err = messageArguments.Marshal()

package message

import "encoding/json"

func UnmarshalMessageArguments(data []byte) (MessageArguments, error) {
	var r MessageArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MessageArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MessageArguments struct {
	Message    string `json:"Message"`    // Message to send
	WebhookURL string `json:"Webhook_URL"`// Slack webhook url
}

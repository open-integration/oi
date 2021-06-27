// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    service, err := UnmarshalService(bytes)
//    bytes, err = service.Marshal()

package types

import "encoding/json"

func UnmarshalService(data []byte) (Service, error) {
	var r Service
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Service) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Service struct {
	Endpoints *Endpoints `json:"endpoints,omitempty"`
}

type Endpoints struct {
	Message *Message `json:"message,omitempty"`
}

type Message struct {
	Arguments *Arguments             `json:"arguments,omitempty"`
	Returns   map[string]interface{} `json:"returns,omitempty"`
}

type Arguments struct {
	Message    *string `json:"Message,omitempty"`     // Message to send
	WebhookURL *string `json:"Webhook_URL,omitempty"` // Slack webhook url
}

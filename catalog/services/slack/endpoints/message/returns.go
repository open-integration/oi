// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    messageReturns, err := UnmarshalMessageReturns(bytes)
//    bytes, err = messageReturns.Marshal()

package message

import "encoding/json"

type MessageReturns map[string]interface{}

func UnmarshalMessageReturns(data []byte) (MessageReturns, error) {
	var r MessageReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MessageReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

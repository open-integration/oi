// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getCardsArguments, err := UnmarshalGetCardsArguments(bytes)
//    bytes, err = getCardsArguments.Marshal()
//
//    getCardsReturns, err := UnmarshalGetCardsReturns(bytes)
//    bytes, err = getCardsReturns.Marshal()
//
//    archiveCardArguments, err := UnmarshalArchiveCardArguments(bytes)
//    bytes, err = archiveCardArguments.Marshal()
//
//    archiveCardReturns, err := UnmarshalArchiveCardReturns(bytes)
//    bytes, err = archiveCardReturns.Marshal()
//
//    addCardArguments, err := UnmarshalAddCardArguments(bytes)
//    bytes, err = addCardArguments.Marshal()
//
//    addCardReturns, err := UnmarshalAddCardReturns(bytes)
//    bytes, err = addCardReturns.Marshal()

package types

import "encoding/json"

func UnmarshalGetCardsArguments(data []byte) (GetCardsArguments, error) {
	var r GetCardsArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetCardsArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetCardsReturns []Card

func UnmarshalGetCardsReturns(data []byte) (GetCardsReturns, error) {
	var r GetCardsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetCardsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalArchiveCardArguments(data []byte) (ArchiveCardArguments, error) {
	var r ArchiveCardArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ArchiveCardArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ArchiveCardReturns map[string]interface{}

func UnmarshalArchiveCardReturns(data []byte) (ArchiveCardReturns, error) {
	var r ArchiveCardReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ArchiveCardReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAddCardArguments(data []byte) (AddCardArguments, error) {
	var r AddCardArguments
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddCardArguments) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AddCardReturns map[string]interface{}

func UnmarshalAddCardReturns(data []byte) (AddCardReturns, error) {
	var r AddCardReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddCardReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetCardsArguments struct {
	Auth  Auth   `json:"Auth"`
	Board string `json:"Board"` // Trello Board ID
}

type Auth struct {
	App   string `json:"App"`   // Trello Application ID
	Token string `json:"Token"` // Trello token
}

type Card struct {
	Badges                Badges        `json:"badges"`
	Board                 interface{}   `json:"Board"`
	Closed                bool          `json:"closed"`
	DateLastActivity      string        `json:"dateLastActivity"`
	Desc                  string        `json:"desc"`
	Due                   *string       `json:"due"`
	DueComplete           bool          `json:"dueComplete"`
	Email                 string        `json:"email"`
	ID                    string        `json:"id"`
	IDAttachmentCover     string        `json:"idAttachmentCover"`
	IDBoard               string        `json:"idBoard"`
	IDCheckLists          []interface{} `json:"idCheckLists"`
	IDLabels              []string      `json:"idLabels"`
	IDList                string        `json:"idList"`
	IDShort               float64       `json:"idShort"`
	Labels                []Label       `json:"labels"`
	List                  List          `json:"List"`
	ManualCoverAttachment bool          `json:"manualCoverAttachment"`
	Name                  string        `json:"name"`
	Pos                   float64       `json:"pos"`
	ShortLink             string        `json:"shortLink"`
	ShortURL              string        `json:"shortUrl"`
	Subscribed            bool          `json:"subscribed"`
	URL                   string        `json:"url"`
}

type Badges struct {
	Attachments        float64 `json:"attachments"`
	CheckItems         float64 `json:"checkItems"`
	CheckItemsChecked  float64 `json:"checkItemsChecked"`
	Comments           float64 `json:"comments"`
	Description        bool    `json:"description"`
	Subscribed         bool    `json:"subscribed"`
	ViewingMemberVoted bool    `json:"viewingMemberVoted"`
	Votes              float64 `json:"votes"`
}

type Label struct {
	Color   string  `json:"color"`
	ID      string  `json:"id"`
	IDBoard string  `json:"idBoard"`
	Name    string  `json:"name"`
	Uses    float64 `json:"uses"`
}

type List struct {
	Closed     bool    `json:"closed"`
	ID         string  `json:"id"`
	IDBoard    string  `json:"idBoard"`
	Name       string  `json:"name"`
	Pos        float64 `json:"pos"`
	Subscribed bool    `json:"subscribed"`
}

type ArchiveCardArguments struct {
	Auth    Auth     `json:"Auth"`
	CardIDs []string `json:"CardIDs,omitempty"` // IDs to archive
}

type AddCardArguments struct {
	Auth        Auth     `json:"Auth"`
	Board       string   `json:"Board"`                 // Trello board ID
	Description *string  `json:"Description,omitempty"` // Trello description to set on card
	Labels      []string `json:"Labels,omitempty"`      // Trello labels to apply on card
	List        string   `json:"List"`                  // Trello list ID
	Name        string   `json:"Name"`                  // Trello card name
}

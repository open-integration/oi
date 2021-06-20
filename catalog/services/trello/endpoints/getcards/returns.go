// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getcardsReturns, err := UnmarshalGetcardsReturns(bytes)
//    bytes, err = getcardsReturns.Marshal()

package getcards

import "encoding/json"

type GetcardsReturns []Card

func UnmarshalGetcardsReturns(data []byte) (GetcardsReturns, error) {
	var r GetcardsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetcardsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
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

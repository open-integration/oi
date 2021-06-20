// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getRecordsReturns, err := UnmarshalGetRecordsReturns(bytes)
//    bytes, err = getRecordsReturns.Marshal()

package getrecords

import "bytes"
import "errors"
import "encoding/json"

func UnmarshalGetRecordsReturns(data []byte) (GetRecordsReturns, error) {
	var r GetRecordsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetRecordsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetRecordsReturnsClass struct {
	Records []Record `json:"Records,omitempty"`
}

type Record struct {
	CreatedTime *string                `json:"createdTime,omitempty"`
	Deleted     *bool                  `json:"deleted,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
	ID          *string                `json:"id,omitempty"`
	Typecast    *bool                  `json:"typecast,omitempty"`
}

type GetRecordsReturns struct {
	AnythingArray          []interface{}
	Bool                   *bool
	Double                 *float64
	GetRecordsReturnsClass *GetRecordsReturnsClass
	Integer                *int64
	String                 *string
}

func (x *GetRecordsReturns) UnmarshalJSON(data []byte) error {
	x.AnythingArray = nil
	x.GetRecordsReturnsClass = nil
	var c GetRecordsReturnsClass
	object, err := unmarshalUnion(data, &x.Integer, &x.Double, &x.Bool, &x.String, true, &x.AnythingArray, true, &c, false, nil, false, nil, true)
	if err != nil {
		return err
	}
	if object {
		x.GetRecordsReturnsClass = &c
	}
	return nil
}

func (x *GetRecordsReturns) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, x.Double, x.Bool, x.String, x.AnythingArray != nil, x.AnythingArray, x.GetRecordsReturnsClass != nil, x.GetRecordsReturnsClass, false, nil, false, nil, true)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}

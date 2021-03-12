package mongolang

/*
	Code related to parsing json to produce bson.D and bson.A structs.

*/

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// JSONToBSON supports parsing JSON strings
// into BSON structures, either bson.D or bson.A
type JSONToBSON struct {
	decoder *json.Decoder
	token   json.Token

	IsBSOND bool
	Err     error

	BSOND bson.D
	BSONA bson.A
}

// ParseJSON parses JSON strings
// to create a bson.D or bson.A struct
func (j *JSONToBSON) ParseJSON(jsonStr string) (interface{}, error) {

	j.decoder = json.NewDecoder(strings.NewReader(jsonStr))
	j.token, j.Err = j.decoder.Token()

	if j.Err != nil {
		return nil, j.Err
	}

	switch v := j.token.(type) {
	case json.Delim:
		if v.String() == "{" {
			j.IsBSOND = true
			j.BSOND = bson.D{}

			var bsonD bson.D
			bsonD, j.Err = j.parseNameValue()

			if j.Err == nil {
				j.BSOND = bsonD
			}
		} else {
			j.IsBSOND = false
			j.BSONA, j.Err = j.parseArray()

		}
	default:
		j.Err = fmt.Errorf("unrecognized or invalid type %T value %v", j.token, j.token)
		return nil, j.Err
	}

	if j.IsBSOND {
		return j.BSOND, j.Err
	}

	return j.BSONA, j.Err

}

// parse a bson.E name value pair
func (j *JSONToBSON) parseNameValue() (bson.D, error) {
	result := bson.D{}

	for {
		j.token, j.Err = j.decoder.Token()
		if j.Err != nil {
			return bson.D{}, j.Err
		}

		switch v := j.token.(type) {
		case string:
			var value interface{}
			value, j.Err = j.parseValue()
			element := bson.E{Key: string(v), Value: value}
			result = append(result, element)
		case json.Delim:
			delim := string(v)
			if delim == "{" {
				fmt.Println("found '{' in parseNameValue()")
				var value interface{}
				value, j.Err = j.parseValue()
				element := bson.E{Key: string(v), Value: value}
				result = append(result, element)
			} else if delim == "}" {
				return result, j.Err
			} else {
				// didn't find the expected "}" delimiter
				j.Err = fmt.Errorf("looking key:value pair, found type %T value %v", j.token, j.token)
				return result, j.Err
			}
		default:
			j.Err = fmt.Errorf("looking key:value pair, found type %T value %v", j.token, j.token)
			return result, j.Err
		}
	}
	// return result, j.Err
}

// parse a json value.
// Either a string, float64, bool, nil, bson.E or bson.A
func (j *JSONToBSON) parseValue() (interface{}, error) {
	j.token, j.Err = j.decoder.Token()
	if j.Err != nil {
		return j.token, j.Err
	}

	switch v := j.token.(type) {
	case string:
		return string(v), nil
	case float64:
		return float64(v), nil
	case bool:
		return bool(v), nil
	case nil:
		return nil, nil
	case json.Delim:
		delim := string(v)
		if delim == "[" {
			return j.parseArray()
		}
		if delim == "{" {
			return j.parseNameValue()
		}

	}

	j.Err = fmt.Errorf("looking for a json value or bson.E or bson.A, found %T value %v", j.token, j.token)
	return nil, j.Err
}

// parse an array of values
func (j *JSONToBSON) parseArray() (bson.A, error) {
	result := bson.A{}

	for {
		j.token, j.Err = j.decoder.Token()
		if j.Err != nil {
			return result, j.Err
		}

		switch v := j.token.(type) {
		case string:
			result = append(result, string(v))
		case float64:
			result = append(result, float64(v))
		case bool:
			result = append(result, bool(v))
		case nil:
			result = append(result, nil)

		case json.Delim:
			delim := string(v)
			if delim == "[" {
				var element bson.A
				element, j.Err = j.parseArray()
				if j.Err != nil {
					return result, j.Err
				}

				result = append(result, element)

			} else if delim == "{" {
				var element bson.D
				element, j.Err = j.parseNameValue()
				if j.Err != nil {
					return result, j.Err
				}
				result = append(result, element)

			} else if delim == "]" {
				return result, j.Err
			}

			// we assume json.Decoder will balance proper {} or []
			// therefore, we have an ending delimiter
			// return the results.

			// return result, j.Err
		}

	}
}

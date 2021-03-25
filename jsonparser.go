package mongolang

/*
	Code related to parsing json to produce bson.D and bson.A structs.

*/

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// JSONToBSON supports parsing JSON strings
// into BSON structures, either bson.D or bson.A
type JSONToBSON struct {
	IsBSOND bool
	IsBSONA bool
	IsOther bool

	Err error

	BSOND bson.D
	BSONA bson.A
	Other interface{}
}

// LastJSONToBSON contains the last instance of this struct
// created in ParseJSONToBSON. This is in case caller wants to
// check the Err or other values before using the parser results.
var LastJSONToBSON JSONToBSON

// ParseJSONToBSON is a convenience function that creates a new
// JSONToBSON struct then uses it to convert the input JSON string
// to a bson.A or bson.D struct.
// Warning... any errors are ignored, but...
// if there is an error, the error is returned. This should cause a later
// error if used in place of a bson struct.
// In a single user TEST environment, caller can also check LastJSONToBSON
func ParseJSONToBSON(jsonStr string) interface{} {
	s := JSONToBSON{}
	result, err := s.ParseJSON(jsonStr)
	LastJSONToBSON = s // for caller error checking

	if err != nil {
		fmt.Printf("error in ParseJSONToBSON: %v \n", err)
		return nil
	}
	return result
}

// ParseJSON parses a JSON string
// to create a bson.D, bson.A or "other" structure.
// Uses the MongoDB canonical (strict) extended JSON format documented at
// https://docs.mongodb.com/manual/reference/mongodb-extended-json/
//
// For example, a filter string searching for specific _id would be coded as
//
//    filter := `{"_id" : {"$oid":"5bf36072a5820f6e28a4736c"} }`
//
// Note use of backticks, NOT single quotes, to create a Go "raw string" which
// can contain unescaped quote marks.
//
// See  https://docs.mongodb.com/manual/reference/mongodb-extended-json/ for
// a complete list of extensions along with examples.
//
// In addition to returning a bson or other struct, ParseJSON sets a series of flags
// to define if the result was a bson.A, bson.D or "other" type. It also sets the
// corresponding type value, BSOND, BSONA, or Other to pointers to the results.
func (j *JSONToBSON) ParseJSON(jsonStr string) (interface{}, error) {

	var doc interface{}
	j.Err = bson.UnmarshalExtJSON([]byte(jsonStr), true, &doc)

	if j.Err != nil {
		return nil, j.Err
	}

	switch v := doc.(type) {
	case bson.D:
		j.BSOND = v
		j.IsBSOND = true
	case bson.A:
		j.BSONA = v
		j.IsBSONA = true
	default:
		j.Other = v
		j.IsOther = true

	}

	return doc, nil
}

package mongolang

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

/*
	Miscellaneous methods that are not related to a struct
*/

// PrintStruct prints an interface object
// as "pretty" json
func PrintStruct(s interface{}) {
	json, _ := json.MarshalIndent(s, "", "  ")
	fmt.Printf("%s\n", json)
}

// PrintBSOND prints a single bson document.
//
// Default is to print JSON in a pretty format
// (one field per line). This can be overridden by passing
// a "pretty" parameter value of false
func PrintBSOND(doc *bson.D, pretty ...bool) {
	prettyPrint := true
	if len(pretty) > 0 {
		prettyPrint = pretty[0]
	}
	if prettyPrint {
		fmt.Printf("\n{ \n")
		for _, v := range *doc {
			fmt.Printf("    %s : %v \n", v.Key, v.Value)
		}
		fmt.Printf("} \n\n")
	} else {
		fmt.Printf("%v \n", doc)
	}
}

// PrintBSONM Prints a bson.M document
func PrintBSONM(doc *bson.M, pretty ...bool) {
	prettyPrint := true
	if len(pretty) > 0 {
		prettyPrint = pretty[0]
	}
	if prettyPrint {
		fmt.Printf("{ \n")
		for k, v := range *doc {
			fmt.Printf("    %s : %v \n", k, v)
		}
		fmt.Printf("} \n")
	} else {
		fmt.Printf("%#v \n", doc)
	}
}

// Allowed Types Flags
// Used to build a uint32 passed to verifyParm.
// Example, to verify that parm is bson.D or bson.M:
//		verifyParm(parm,(bsonDAllowed|bsonMAllowed))
const (
	bsonDAllowed      = 1 // bit 1
	bsonMAllowed      = 2 // bit 2
	bsonAAllowed      = 4 // bit 3
	bsonDSliceAllowed = 8 // bit 4
)

// Verify the type of a parameter based on allowedTypesFlags.
// In two cases will convert from one type to another:
//   1. If parm is a bson.D and bson.D is not allowed
//		but a bson.A is allowed,
//      will return the bson.D wrapped in a bson.A or []bson.D
//   2. If nil passed, will return an empty bson.D, bson.A, or bson.M
//      if one of those is allowed.
func verifyParm(parm interface{}, allowedTypes uint32) (interface{}, error) {

	switch p := parm.(type) {
	case string: // parse strings
		s := JSONToBSON{}
		result, err := s.ParseJSON(p)

		if err != nil {
			fmt.Printf("error in ParseJSONToBSON: %v \n", err)
			return nil, err
		}

		parm = result
	}

	switch parm.(type) {
	case nil:
		if allowedTypes&bsonDAllowed != 0 {
			return bson.D{}, nil
		}
		if allowedTypes&bsonAAllowed != 0 {
			return bson.A{}, nil
		}
		if allowedTypes&bsonMAllowed != 0 {
			return bson.M{}, nil
		}
		return parm, errors.New("nil parm without suitable default type")

	case bson.D:
		if allowedTypes&bsonDAllowed != 0 {
			return parm, nil
		}
		if allowedTypes&bsonAAllowed != 0 {
			return bson.A{parm}, nil
		}

	case bson.M:
		if allowedTypes&bsonMAllowed != 0 {
			return parm, nil
		}

	case bson.A:
		if allowedTypes&bsonAAllowed != 0 {
			return parm, nil
		}

	case []bson.D:
		if allowedTypes&bsonDSliceAllowed != 0 {
			return parm, nil
		}
	}

	return nil, fmt.Errorf("invalid parm type: %T", parm)
}

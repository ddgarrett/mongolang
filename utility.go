package mongolang

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

/*
	Miscellaneous methods that are not related to a struct
*/

// convertBSONParm converts an entry in an array of interface{}
// to either a bson.M or bson.D interface
//
// i = index to array of parms
// parms = variable number of parms which can be bson.M or bson.D interfaces
//
// bson.M{} interface is returned if parm is not a bson.M or bson.D
/*  TODO: modify logic to use simpler logic to convert  an interface{} such as:
switch v := i.(type) {
case int:
    fmt.Println("twice i is", v*2)
case float64:
    fmt.Println("the reciprocal of i is", 1/v)
case string:
    h := len(v) / 2
    fmt.Println("i swapped by halves is", v[h:]+v[:h])
default:
    // i isn't one of the types above
}

*/
func convertBSONParm(i int, parms ...interface{}) interface{} {
	if len(parms) > i && parms[i] != nil {
		parm := parms[i].([]interface{})

		if len(parm) == 0 {
			return bson.M{}
		}

		reflectType := reflect.TypeOf(parm[0])

		if reflectType.PkgPath() == "go.mongodb.org/mongo-driver/bson/primitive" {
			if reflectType.Name() == "M" {
				return parm[0].(bson.M)
			} else if reflectType.Name() == "D" {
				return parm[0].(bson.D)
			}
		}
	}
	return bson.M{}
}

// PrintStruct prints an interface object
// as "pretty" json
func PrintStruct(s interface{}) {
	json, _ := json.MarshalIndent(s, "", "  ")
	fmt.Printf("data:\n%s\n", json)
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
		fmt.Printf("\n{ \n")
		for k, v := range *doc {
			fmt.Printf("    %s : %v \n", k, v)
		}
		fmt.Printf("} \n\n")
	} else {
		fmt.Printf("%#v \n", doc)
	}
}

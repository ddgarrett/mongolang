package mongolang

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
// In some cases will convert from one type to another:
//   1. If parm is a bson.D and bson.D is not allowed
//		but a bson.A is allowed,
//      will return the bson.D wrapped in a bson.A or []bson.D
//   2. If nil passed, will return an empty bson.D, bson.A, or bson.M
//      if one of those is allowed.
//   3. If a string is passed, will assume it is an extended JSON string
//      and call bson.UnmarshalExtJSON to parse the JSON string.
func verifyParm(parm interface{}, allowedTypes uint32) (interface{}, error) {

	switch p := parm.(type) {
	case string: // parse strings
		var result interface{}

		err := bson.UnmarshalExtJSON([]byte(p), true, &result)

		if err != nil {
			fmt.Printf("error in ParseJSONToBSON: %v \n", err)
			return nil, err
		}

		parm = result
	}

	switch p := parm.(type) {
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

		if allowedTypes&bsonDSliceAllowed != 0 {
			invalidSubtype := false
			r2 := make([]bson.D, 0, len(p))
			for _, v := range p {
				v2, ok := v.(bson.D)
				if !ok {
					invalidSubtype = true
				} else {
					r2 = append(r2, v2)
				}
			}

			if !invalidSubtype {
				return r2, nil
			}
		}

	case []bson.D:
		if allowedTypes&bsonDSliceAllowed != 0 {
			return parm, nil
		}
	}

	return nil, fmt.Errorf("invalid parm type: %T", parm)
}

type printBSONParms struct {
	indent      int
	prevBracket bool
}

// PrintBSON prints formatted BSON structures,
// optionally with the value type.
func PrintBSON(parm interface{}) {

	parms := printBSONParms{0, false}

	parms.printBSON(parm)

}

// NOTE: returns true if parm was a BSON type
func (p *printBSONParms) printBSON(parm interface{}) (isBSON bool) {

	switch pt := parm.(type) {
	case bson.E:
		p.printBSONE(pt)
	case bson.D:
		p.printBSOND(pt)
	case bson.A:
		p.printBSONA(pt)
	default:
		s := fmt.Sprintf("%v", pt)
		if len(s) > 40 {
			s = s[:30] + "..."
		}
		fmt.Printf("%v", s)
		return false
	}

	return true
}

func (p *printBSONParms) printIndent() {
	fmt.Println()
	fmt.Printf("%*s", p.indent*3, " ")
}

func (p *printBSONParms) printBSONE(parm bson.E) {

	p.prevBracket = false

	tString := fmt.Sprintf("%T", parm.Value)
	tString = strings.TrimPrefix(tString, "primitive.")

	p.printIndent()

	fmt.Printf("%s (%s): ", parm.Key, tString)

	p.printBSON(parm.Value)
	p.prevBracket = false
}

func (p *printBSONParms) printBSOND(parm bson.D) {

	if p.prevBracket {
		fmt.Printf(",")
		p.printIndent()
	}

	fmt.Printf("{")
	p.indent++
	p.prevBracket = false

	for _, v := range parm {
		p.printBSON(v)
	}

	p.indent--
	p.printIndent()
	fmt.Printf("}")
	p.prevBracket = true
}

func (p *printBSONParms) printBSONA(parm bson.A) {
	// printIndent(indent)
	fmt.Printf("[")
	p.indent++
	p.prevBracket = false

	// An array of BSON objects is printed one object per line
	// but an array of elementary items is printed on a single
	// line, separated by ', '
	isBSON := false
	lastParmIndex := len(parm) - 1
	for i, v := range parm {
		isBSON = p.printBSON(v)
		if !isBSON && i != lastParmIndex {
			fmt.Printf(", ")
		}
	}

	p.indent--

	if isBSON {
		p.printIndent()
	}
	fmt.Printf("]")
	p.prevBracket = false
}

package mongolang

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func ExamplePrintStruct() {
	// bsonD := bson.D{{Key: "bsonDKey", Value: "bsonDValue"}}
	bsonD := bson.D{{Key: "bsonDKey", Value: "bsonDValue"}}
	PrintStruct(bsonD)
	// output:
	// [
	//   {
	//     "Key": "bsonDKey",
	//     "Value": "bsonDValue"
	//   }
	// ]

}

func TestVerifyParmsPart01(t *testing.T) {
	bsonD := bson.D{{Key: "bsonDKey", Value: "bsonDValue"}}
	bsonM := bson.M{"bsonMKey": "bsonMValue"}
	// bsonA := bson.A{bsonD, bsonM}

	const allowAD = bsonDAllowed | bsonAAllowed
	// const allowDM = bsonDAllowed | bsonMAllowed

	// Test that verify bson.D works
	parm, err := verifyParm(bsonD, allowAD)
	_, okType := parm.(bson.D)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParms t01 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that bson.M will be rejected
	parm, err = verifyParm(bsonM, allowAD)
	if err == nil {
		t.Errorf("TestVerifyParms t02 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that bson.D will be wrapped in a bson.M if needed
	parm, err = verifyParm(bsonD, bsonAAllowed)
	_, okType = parm.(bson.A)
	if !okType || parm == nil || err != nil {
		t.Errorf("TestVerifyParms t03 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that JSON string bson.D works
	stringD := `{"bsonDKey2":"bsonDValue2"}`
	parm, err = verifyParm(stringD, bsonDAllowed)
	_, okType = parm.(bson.D)
	if !okType || parm == nil || err != nil {
		t.Errorf("TestVerifyParms t04 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that JSON string bson.A works
	stringA := `[{"anotherKey":"anotherValue"},{"key3":"value3"}]`
	parm, err = verifyParm(stringA, bsonAAllowed)
	_, okType = parm.(bson.A)
	if !okType || parm == nil || err != nil {
		t.Errorf("TestVerifyParms t05 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// and that the element in bson.A is a bson.D
	pa, okType := parm.(bson.A)
	if okType {
		parm2 := pa[0]
		_, okType := parm2.(bson.D)
		if !okType {
			t.Errorf("TestVerifyParms t06 parm: %v, type: %T", parm2, parm2)
		}
	}

	// Test that invalid JSON string is caught
	invalidString := `["missingBracket":"noBracketValue"]`
	parm, err = verifyParm(invalidString, allowAD)
	if parm != nil || err == nil {
		t.Errorf("TestVerifyParms t07 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that verify bson.M works
	parm, err = verifyParm(bsonM, bsonMAllowed)
	_, okType = parm.(bson.M)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParms t089 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that verify bson.D slice works
	bsonDSlice := []bson.D{}
	parm, err = verifyParm(bsonDSlice, bsonDSliceAllowed)
	_, okType = parm.([]bson.D)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParms t089 parm: %v, type: %T, error: %v", parm, parm, err)
	}
}

func TestVerifyParmsNil(t *testing.T) {
	// Test cases where input parm is nil

	var nilParm interface{} = nil

	// Test that verify bson.D works
	parm, err := verifyParm(nilParm, bsonDAllowed)
	_, okType := parm.(bson.D)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParmsNil t01 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that verify bson.M works
	parm, err = verifyParm(nilParm, bsonAAllowed)
	_, okType = parm.(bson.A)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParmsNil t02 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that verify bson.A works
	parm, err = verifyParm(nilParm, bsonMAllowed)
	_, okType = parm.(bson.M)
	if parm == nil || err != nil || !okType {
		t.Errorf("TestVerifyParmsNil t03 parm: %v, type: %T, error: %v", parm, parm, err)
	}

	// Test that invalid nil parm is caught
	parm, err = verifyParm(nilParm, 0)
	if parm != nil || err == nil {
		t.Errorf("TestVerifyParmsNil t04 parm: %v, type: %T, error: %v", parm, parm, err)
	}
}

func ExamplePrintBSON() {
	pipeline := `[
        { "$match" : {"_id" : {"$oid":"5bf36072a5820f6e28a4736c"} }},
		{ "$test" : "field that will be truncated because of length" },
		{ "$test" : ["an array", "of text", "values"]},
		{ "$test" : [{"array": "of"}, {"objects": 2}]},
		{ "$limit": 3 }
	]`

	var doc interface{}
	err := bson.UnmarshalExtJSON([]byte(pipeline), true, &doc)
	if err != nil {
		panic(err)
	}

	PrintBSON(doc)

	// output:
	// [{
	//       $match (D): {
	//          _id (ObjectID): ObjectID("5bf36072a5820f6e28a4736c")
	//       }
	//    },
	//    {
	//       $test (string): field that will be truncated b...
	//    },
	//    {
	//       $test (A): [an array, of text, values]
	//    },
	//    {
	//       $test (A): [{
	//             array (string): of
	//          },
	//          {
	//             objects (int32): 2
	//          }
	//       ]
	//    },
	//    {
	//       $limit (int32): 3
	//    }
	//  ]
}

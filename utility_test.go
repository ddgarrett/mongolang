package mongolang

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestVerifyParms(t *testing.T) {
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

}

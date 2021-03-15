package mongolang

import (
	"strings"
	"testing"
)

func TestSkip(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")

	result := db.Coll("zips").
		Find(`{"state":"CA"}`, `{"loc":0}`).Sort(`{"pop":-1}`).Skip(2).
		Limit(1).Pretty()

	// fmt.Printf("%s", result)
	i := strings.Index(result, "_id : 90650")

	if i < 0 {
		t.Errorf("TestSkip didn't find zipcode 90650: \n%s\n", result)
	}
}

func TestCount(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")

	count := db.Coll("zips").
		Find(`{"state":"CA"}`).Count()

	if count != 1516 {
		t.Errorf("TestCount had %d instead of 1516 as count", count)
	}

}

func TestString(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")

	result := db.Coll("zips").
		Find(`{"state":"CA"}`, `{"loc":0}`).Sort(`{"pop":-1}`).Skip(2).
		Limit(1).String()

	// fmt.Printf("%s", result)
	i := strings.Index(result, `"Value": "NORWALK"`)

	if i < 0 {
		t.Errorf("TestSkip didn't find \"Value\": \"NORWALK\" in \n%s", result)
	}
}

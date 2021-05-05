package mongolang

import (
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestSkip(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

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
	defer db.Disconnect()

	count := db.Coll("zips").
		Find(`{"state":"CA"}`).Count()

	if count != 1516 {
		t.Errorf("TestCount had %d instead of 1516 as count", count)
	}

}

func TestString(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	result := db.Coll("zips").
		Find(`{"state":"CA"}`, `{"loc":0}`).Sort(`{"pop":-1}`).Skip(2).
		Limit(1).String()

	// fmt.Printf("%s", result)
	i := strings.Index(result, `"Value": "NORWALK"`)

	if i < 0 {
		t.Errorf("TestSkip didn't find \"Value\": \"NORWALK\" in \n%s", result)
	}
}

// TestErr tests various cursor error conditions
func TestCursorErr(t *testing.T) {

	cursor := &Cursor{}
	if cursor.Err() == nil {
		t.Error("Err() did not return invalid cursor error")
	}

	if cursor.requireOpenCursor() {
		t.Error("requireOpenCursor() returned true for invalid cursor")
	}

	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017")
	defer db.Disconnect()

	// Make sure error returned is not connected to db
	cursor = db.Coll("zips").Find().Limit(1)
	if cursor.Err() != ErrNotConnectedDB {
		t.Errorf("expected ErrNotConnectedDB, got %v", cursor.Err())
	}

	// Test closed cursor
	db.Use("quickstart")
	cursor = db.Coll("zips").Find(`{"state":"CA"}`)
	cursor.HasNext() // make sure we opened a  Mongo cursor
	cursor.Close()

	// try to close twice
	e := cursor.Close()
	if e != ErrClosedCursor {
		t.Errorf("expected ErrClosedCursor from Close() call got %v", e)
	}

	// other calls on closed cursor
	cursor.Sort(`{"zip":1}`).HasNext()
	cursor.Pretty()
	cursor.ToArray()
	cursor.Next()
	cursor.Count()
	s := cursor.String()
	if cursor.Err() != ErrClosedCursor {
		t.Errorf("expected ErrClosedCursor got %v, String(): %s", cursor.Err(), s)
	}

	// invalid sort parms
	db.Coll("zips").Find(`{"state":"CA"}`).Sort(bson.A{})
	if db.Err == nil {
		t.Error("expected error from invalid sort parm")
	}

	// test HasNext() with buffered doc
	cursor = db.Coll("zips").Find(`{"state":"CA"}`)
	cursor.HasNext()
	cursor.HasNext()

	// next with no next available
	cursor = db.Coll("zips").Find(`{"state":"abc"}`)
	cursor.Next()
	if cursor.Err() == nil {
		t.Error("expected error from invalid Next() call")
	}
}

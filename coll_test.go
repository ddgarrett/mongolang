package mongolang

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/message"
)

type cityState struct {
	State string
	City  string
}

type cityStatePop struct {
	ID  cityState `bson:"_id"`
	Pop int
}

func ExampleAggregate() {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	// Three largest cities
	pipeline := `[
		{ "$group":
			{
				"_id": { "State": "$state", "City": "$city" },
				"Pop": { "$sum": "$pop" }
			}
		},
		{"$sort": {"Pop":-1}},
		{ "$limit": 3 }
	]`

	var r []cityStatePop
	db.Coll("zips").Aggregate(pipeline).ToArray(&r)

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("result has %d cities. Largest is: %s, %s with population of %d\n",
		len(r), r[0].ID.City, r[0].ID.State, r[0].Pop)

	// output: result has 3 cities. Largest is: CHICAGO, IL with population of 2,452,177
}

func ExampleFind() {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	cursor := db.Coll("zips").
		Find(`{"state":"CA", "pop":{"$gt":1000}}`, `{"loc":0}`).
		Sort(`{"pop":-1}`).Limit(3)

	if cursor.HasNext() {
		fmt.Println("largest: ", cursor.Next())
	}

	fmt.Println(cursor.Next())

	var r []bson.M
	cursor.ToArray(&r)

	rLargest := r[0]

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("result has %d zipcodes. Largest zip code in CA: %s, %s with population of %d\n",
		len(r), rLargest["city"], rLargest["_id"], rLargest["pop"])

	// output:
	// largest:  &[{_id 90201} {city BELL GARDENS} {pop 99568} {state CA}]
	// &[{_id 90011} {city LOS ANGELES} {pop 96074} {state CA}]
	// result has 3 zipcodes. Largest zip code in CA: BELL GARDENS, 90201 with population of 99,568

}

func ExampleFindOne() {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	result := db.Coll("zips").
		FindOne(`{"_id":"90002"}`, `{"loc":0}`)

	fmt.Println("CA zip: ", result)

	// output:
	// CA zip:  &[{_id 90002} {city LOS ANGELES} {pop 40629} {state CA}]
}

func TestFindOne(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()
	result := db.Coll("zips").FindOne()

	resultE := ([]primitive.E)(*result)
	if len(resultE) != 5 {
		t.Errorf("TestFindOne expected to find a doc with 5 fields instead of %d field(s)", len(resultE))
	}
}
func TestInsertOne(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	// test invalid insert document type
	db.Coll("testCollection").InsertOne(bson.A{})

	if db.Err == nil {
		t.Error("TestInsertOne with invalid document type expected error")
	}

	// test insert of valid document
	insertDocJSON := `{
		"title": "The Polyglot Developer Podcast",
		"author": "Nic Raboy",
		"tags": ["development", "programming", "coding"] }`

	result := db.Coll("testCollection").InsertOne(insertDocJSON)

	if db.Err != nil {
		t.Errorf("TestInsertOne insert error %v", db.Err)
	}

	searchKey := bson.M{"_id": result.InsertedID}
	insertedDoc := []bson.M{}
	db.Coll("testCollection").Find(searchKey).ToArray(&insertedDoc)

	if len(insertedDoc) == 0 {
		t.Errorf("TestInsertOne unable to read inserted doc with %v", searchKey)
	}

	title, foundAuthor := insertedDoc[0]["title"]

	if !foundAuthor || title != "The Polyglot Developer Podcast" {
		t.Errorf("TestInsertOne didn't find inserted title The Polyglot Developer Podcast")
	}

	deleteResult := db.Coll("testCollection").DeleteOne(searchKey)

	if db.Err != nil {
		t.Errorf("TestInsertOne delete error %v", db.Err)
	}

	if deleteResult.DeletedCount != 1 {
		t.Errorf("TestInsertOne delete count of %d", deleteResult.DeletedCount)
	}
}

func TestInsertMany(t *testing.T) {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")
	defer db.Disconnect()

	// test insert of invalid document
	db.Coll("testCollection").InsertMany(`a`)

	if db.Err == nil {
		t.Error("TestInsertMany insert of invalid JSON expected error")
	}

	// just in case, delete any existing documents
	db.Coll("testCollection").DeleteMany(`{}`)

	insertDocJSON := `[
		{ "title": "The Polyglot Developer Podcast",
		  "author": "Nic Raboy",
		  "tags": ["development", "programming", "coding"] },

		{ "title": "The Polyglot Developer Podcas Version 2",
			"author": "Nic Raboy",
			"tags": ["development", "programming", "coding"] },

		{ "title": "The Polyglot Developer Podcas Version 3",
			  "author": "Nic Raboy Jr.",
			  "testCase": "Nic Raboy Jr. will not be deleted first time",
			  "tags": ["development", "programming", "coding"] }
	]`

	result := db.Coll("testCollection").InsertMany(insertDocJSON)

	if db.Err != nil {
		t.Errorf("TestInsertMany insert error: %v", db.Err)
		return
	}

	if len(result.InsertedIDs) != 3 {
		t.Errorf("TestInsertMany only inserted %d docs", len(result.InsertedIDs))
	}

	// Delete many for author Nic Raboy should delete 2, leaving 1 behind.
	// Use deleteMany a second time to delete remaining document.

	deleteResult := db.Coll("testCollection").DeleteMany(`{"author":"Nic Raboy"}`)

	if db.Err != nil {
		t.Errorf("TestInsertOne delete error %v", db.Err)
	}

	if deleteResult.DeletedCount != 2 {
		t.Errorf("TestInsertOne DeleteMany count of %d. Expected 2", deleteResult.DeletedCount)
	}

	// delete the remaining inserted document
	db.Coll("testCollection").DeleteMany(`{"author":"Nic Raboy Jr."}`)
}

// dbTest tests that we received the expected error from a call
// where the DB is not connected
func testErrNotConnectedDB(db DB, t *testing.T, f string) {
	if db.Err != ErrNotConnectedDB {
		t.Errorf("expected ErrNotConnectedDB after %s, got: %v", f, db.Err)
	}
}

// TestErr tests various error conditions
func TestErr(t *testing.T) {

	// test for invalid connection
	coll := &Coll{}

	if coll.Err() != ErrInvalidColl {
		t.Errorf("expected invalid collection error, got: %v", coll.Err())
	}

	// test FindOne(), etc. without specifying a DB to use first
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017")
	defer db.Disconnect()

	db.Coll("zips").FindOne()
	testErrNotConnectedDB(db, t, "FindOne()")

	db.Coll("zips").Find()
	testErrNotConnectedDB(db, t, "Find()")

	db.Coll("zips").Aggregate("[]")
	testErrNotConnectedDB(db, t, "Aggregate()")

	db.Coll("testCollection").InsertOne("{}")
	testErrNotConnectedDB(db, t, "InsertOne()")

	db.Coll("testCollection").InsertMany("[]")
	testErrNotConnectedDB(db, t, "InsertMany()")

	db.Coll("testCollection").DeleteOne("{}")
	testErrNotConnectedDB(db, t, "DeleteOne()")

	db.Coll("testCollection").DeleteMany("{}")
	testErrNotConnectedDB(db, t, "DeleteMany()")

	// test for reset error in FindOne()
	db.Use("quickstart")
	coll = db.Coll("zips")
	coll.setErr(ErrNotConnected)
	coll.FindOne()

	if coll.Err() != nil {
		t.Errorf("expected nil error after FindOne, got: %v", coll.Err())
	}

	// test FindOneErrors
	coll.FindOne(bson.A{})
	if coll.Err() == nil {
		t.Error("expected error after FindOne() with bson.A{} filter")
	}

	coll.FindOne(bson.M{}, bson.A{})
	if coll.Err() == nil {
		t.Error("expected error after FindOne() with bson.A{} options")
	}

	coll.FindOne(`{"zipCode":{"$$$invalid":0}}}`)
	if coll.Err() == nil {
		t.Error("expected error after FindOne() with invalid filter")
	}

	// test Find() errors
	coll.Find(bson.A{})
	if coll.Err() == nil {
		t.Error("expected error after Find() with bson.A{} filter")
	}

	coll.Find(bson.M{}, bson.A{})
	if coll.Err() == nil {
		t.Error("expected error after Find() with bson.A{} options")
	}

	// test that find works okay without any parms
	coll.Find()
	if coll.Err() != nil {
		t.Errorf("expected nil error after Find() without parms, got %v", coll.Err())
	}

	// test delete with invalid filter
	db.Coll("testCollection").DeleteOne(bson.A{})
	if db.Err == nil {
		t.Error("DeleteOne with invalid filter JSON expected error")
	}

	db.Coll("testCollection").DeleteMany(`a`)
	if db.Err == nil {
		t.Error("DeleteMany with invalid filter JSON expected error")
	}

}

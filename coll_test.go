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

	author, foundAuthor := insertedDoc[0]["author"]

	if !foundAuthor || author != "Nic Raboy" {
		t.Errorf("TestInsertOne didn't find inserted author of Nic Raboy")
	}

	// TODO: delete inserted document or drop collection
}

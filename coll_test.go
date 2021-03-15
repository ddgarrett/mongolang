package mongolang

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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

func ExampleFindOne() {
	db := DB{}
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")

	cursor := db.Coll("zips").
		Find(`{"state":"CA", "pop":{"$gt":1000}}`, `{"loc":0}`).
		Sort(`{"pop":-1}`).Limit(3)

	fmt.Println("largest: ", cursor.Next())

	var r []bson.M
	cursor.ToArray(&r)

	rLargest := r[0]

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("result has %d zipcodes. Largest zip code in CA: %s, %s with population of %d\n",
		len(r), rLargest["city"], rLargest["_id"], rLargest["pop"])

	// output:
	// largest:  &[{_id 90201} {city BELL GARDENS} {pop 99568} {state CA}]
	// result has 3 zipcodes. Largest zip code in CA: BELL GARDENS, 90201 with population of 99,568

}

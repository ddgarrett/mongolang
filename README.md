# MonGolang
## `MongoDB`
## ` +  Golang`

### Goal

Simplify Go programs which use MongoDB.
Create a framework which supports the 10% of the calls which are needed 90% of the time, while easily allowing the use of the other 90% of MongoDB capabilities which are needed 10% of the time.

1. Support the 10% of the calls which are used 90% of the time
2. The other 90% of the calls which are used 10% of the time should be no more difficult to do then without this framework
3. Calls should resemble the calls made via the MongoDB Console
4. Although error checking is important, it shouldn't get in the way of being able to chain calls such as `db.coll("someCollection").find(bson.M{"lastName":"Johnson"}).sort(bson.A{})limit(10)
6. Don't extend the capabilities of the MongoDB Golang driver, just simplify the use of what's there
7. New MongoDB releases should not require changes to the framework

Accordingly, the following are outside scope:
1. Connection pools
2. Simplification of calls for less used capabilities:
   1. Batch Insert
   2. Create Index
   3. Drop Collection
   4. ... more?
3. 


In Scope:
1. List calls from MongoDB Console
2. xxx


### Problem with Using Golang with MongoDB

1. Use of context, when and how
2. Lots of "cruft" with use of context and error checking
3. Confusion over use of bson.M vs bson.D
  * Use bson.M for everything except insert of documents or other rare instances where order of fields matters, such as sort order?
  * Use bson.D for project, aggregation pipeline steps

4. Are people releasing resources efficiently?

For selection criteria maybe not a big deal?

In JavaScript we can write a selection such as this:


```javascript
  var favorites = ["Sandra Bullock","Tom Hanks","Julia Roberts",
              "Kevin Spacey","George Clooney"]

  var selection = {"countries" : "USA",
                  "tomatoes.viewer.rating" : {"$gte":3},
                  "cast": {"$in":favorites} }

```

In golang this would be:

```go

	favorites := bson.A{"Sandra Bullock", "Tom Hanks", "Julia Roberts",
		"Kevin Spacey", "George Clooney"}

	selection := bson.M{"countries": "USA",
		"tomatoes.viewer.rating": bson.M{"$gte": 3},
		"cast": bson.M{"$in": favorites}}

```

Not too bad with the one caveat that `bson.M` is an *unordered* map. Meaning that the order of the entries in the map are random. This isn't a problem when specifying search criteria and similar parameters. 

The problem comes when specifying new documents. In that case you want to preserve the order of the fields to make the document more human readable. Maybe it doesn't matter to computers if the`"_id"` field is the first, the last or somewhere buried in the middle of a large document, but a human expects the `"_id"` field to be at the beginning of the document where it's easy to find.

When inserting new documents you're probably best off doing it via a `struct` in any case. This does **(?)** preserver the order of the fields? (need to check this out)

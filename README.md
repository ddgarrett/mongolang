# MonGolang
## `MongoDB`
## ` +  Golang`
## `            = Mon Golang`
### Goal

Simplify Go programs which use MongoDB.
Create a framework which supports the 10% of the calls which are needed 90% of the time, while easily allowing the use of the other 90% of MongoDB capabilities which are needed 10% of the time.

1. Support the 10% of the calls which are used 90% of the time
2. The other 90% of the calls which are used 10% of the time should be no more difficult to do then without this framework
3. Calls should resemble the calls made via the MongoDB Console
4. Although error checking is important, it shouldn't get in the way of being able to chain calls such as `db.Coll("someCollection").Find(``{"lastName":"Johnson"}``).Sort('{"lastName":1}).Limit(5).Pretty()`
5. Don't extend the capabilities of the MongoDB Golang driver, just simplify the use of what's there
6. New MongoDB releases should not require changes to the framework

### With Jupyter Notebook
Combined with Jupyter Notebook using GopherNotes, this provides a simple interface to run MongoDB Shell like commands using Go. Not only easy to iterate and easy to rerun, but also a great way to prototype code which can then be copied and pasted into a Go program. 

![Example Using MonGolang in Jupyter Notebook](misc/MonGolang_V01.1._Test01.png?raw=true)


[See here](https://docs.google.com/presentation/d/1zq8-n0w0uiy9AIK9kaOiZgIL6VEmUc1FBDpbImZ4RLw/edit?usp=sharing) for a brief slide presentation on running MonGolang on Jupyter and then using the code in a standalone compiled Go program.

### Some things not too bad?

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





package mongolang

/*
	Methods to support MongoDB Cursor.

	In keeping with golang requirements, all accessible methods start with uppercase.

	Note that a number of calls return a cursor and the actual DB read
	does not actually start until a document is returned from a call.

	This makes the following call possible:

	   cursor := db.Coll("podcosts").Find().Sort(bson.M{}).Skip(10).Limit(100)

	In the above call note that:
	   - .Find() must be the first cursor method - it actually creates the cursor struct
	   - .Sort(), .Skip(), .Limit() can be in any order after the .find()
	   - no documents have yet been returned and therefore no DB reads have occurred

	The following commands will then actually retrieve a document and begin db reads.
	Once these cursor methods are called none of the previously mentioned methods can be called
	on the same cursor.
		- HasNext()		- Next()		- ForEach()
		- ToArray()		- Count()		- Pretty()
		- Close()		- IsClosed

	HasNext() and Next() will access the next document in the cursor (if one exists).
	If no next document exists, the cursor will be closed.
	IsClosed() returns true if the cursor is closed.
	All others will close the cursor after reading all of the documents for the cursor.
	Once the cursor is closed any attempt to use it will cause a panic.

	TODO: methods which read the remaining documents should check to see if one was buffered via
	a HasNext()
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Close closes a cursor
// Note that it's possible to reuse the cursor, though not recommended?
func (c *Cursor) Close() error {
	if !c.IsClosed {
		if c.MongoCursor != nil {
			c.Err = c.MongoCursor.Close(context.Background())
			c.MongoCursor = nil
		}

		c.IsClosed = true

		c.NextDoc = nil

		c.Filter = nil
		c.FindOptions = options.FindOptions{}

		c.AggrPipeline = nil
		c.AggrOptions = options.AggregateOptions{}
	} else {
		c.Err = errors.New("Close called on already closed cursor")
	}

	return c.Err
}

// getMongoCursor ensures that we have an open MongoDB Cursor.
// If the cursor is currently nil, creates a new Find or Aggregate cursor.
func (c *Cursor) getMongoCursor() error {

	if c.MongoCursor == nil {
		var err error
		if c.IsFindCursor {
			c.MongoCursor, err = c.Collection.MongoColl.Find(context.Background(), c.Filter, &c.FindOptions)
		} else {
			c.MongoCursor, err = c.Collection.MongoColl.Aggregate(context.Background(), c.AggrPipeline, &c.AggrOptions)
		}

		// mark as not closed here so that if error, c.Close() reinitializes cursor
		c.IsClosed = false

		if err != nil {
			c.Close()
			c.Err = err
			return err
		}

	}
	return nil
}

// bufferNext reads the next document, if it exists, into the Cursor buffer.
// Returns true if there is a next document.
// If no next document, automatically closes the cursor.
// Since this is an internal method it will panic if called with a closed cursor.
func (c *Cursor) bufferNext() bool {
	if c.IsClosed {
		panic("internal error: bufferNext() called after cursor closed")
	}

	err := c.getMongoCursor()

	if err != nil {
		return false
	}

	hasNext := c.MongoCursor.Next(context.Background())

	if !hasNext {
		c.Close()
		return false
	}

	c.NextDoc = &bson.D{}
	err = c.MongoCursor.Decode(c.NextDoc)

	if err != nil {
		c.Close()
		// Close will set c.Err
		// Don't lose error from Decode call
		c.Err = err
		c.NextDoc = nil
		return false
	}

	c.Err = nil
	return true
}

// sort, skip, limit - pre-cursor open methods
//
// NOT allowed on aggregate cursor

// Sort specifies the bson.D to be used to sort the cursor results
func (c *Cursor) Sort(sortSequence interface{}) *Cursor {

	if c.IsClosed {
		c.Err = errors.New("Sort called on closed cursor")
	} else if !c.IsFindCursor {
		c.Err = errors.New("Sort called on aggregation cursor")
		c.Close()
	} else {
		c.FindOptions.Sort, c.Err = verifyParm(sortSequence, (bsonDAllowed | bsonMAllowed))
	}

	return c
}

// Skip specifies the number of documents to skip before returning the first document
func (c *Cursor) Skip(skipCount int64) *Cursor {
	if c.IsClosed {
		c.Err = errors.New("Skip called on closed cursor")
	} else if !c.IsFindCursor {
		c.Err = errors.New("Skip called on aggregation cursor")
		c.Close()
	} else {
		c.FindOptions.Skip = &skipCount
	}

	return c
}

// Limit specifies the max number of documents to return
func (c *Cursor) Limit(limitCount int64) *Cursor {
	if c.IsClosed {
		c.Err = errors.New("Limit called on closed cursor")
	} else if !c.IsFindCursor {
		c.Err = errors.New("Limit called on aggregation cursor")
		c.Close()
	} else {
		c.Err = nil
		c.FindOptions.Limit = &limitCount
	}

	return c
}

/*
	HasNext() and Next() - used to read through a cursor.
	Closes the cursor if no next document.

*/

// HasNext returns true if the cursor has a next document available.
func (c *Cursor) HasNext() bool {
	if c.IsClosed {
		c.NextDoc = nil
		c.Err = errors.New("HasNext() called on closed cursor")
		return false
	}

	if c.NextDoc != nil {
		return true
	}

	return c.bufferNext()
}

// Next returns the next document for the cursor as a bson.D struct.
// TODO: allow a struct to be passed similar to cursor.ToArray(...)
func (c *Cursor) Next() *bson.D {
	if c.IsClosed {
		c.Err = errors.New("Next() called on closed cursor")
		return &bson.D{}
	}

	if c.NextDoc == nil {
		hasNext := c.bufferNext()
		if !hasNext {
			c.Err = errors.New("Next() called when there isn't a next document")
			return &bson.D{}
		}
	}

	// "unbuffer" the next document
	doc := c.NextDoc
	c.NextDoc = nil

	return doc
}

/*
	Read all of the documents for a cursor then close the cursor.

	- ForEach()		- ToArray()		- Count()		- Pretty()
	- String() (fulfills the Stringer interface for printing, etc.)

*/

// ForEach calls the specified function once for each remaining cursor document
// passing the function a bson.D document.
func (c *Cursor) ForEach(f func(*bson.D)) {
	fmt.Println("ForEach(...) not yet implemented")
}

// ToArray returns all of the remaining documents for a cursor
// in a bson.D slice. NOTE: currently seems to return all docs
// even those already read via a cursor.Next() call.
// Optional parm is a pointer to a slice which typically would contain
// a custom struct or bson.A struct. In this case, ToArray returns an
// empty []bson.D slice.
// This may change at some future date but currently it is difficult
// to deal with a slice of any type. For consistency with Mongo Shell
// ToArray should return a slice. However we also need the ability
// to return all of the documents to a slice which is a custom struct
// or a bson.A struct.
func (c *Cursor) ToArray(parm ...interface{}) []bson.D {
	result := []bson.D{}

	if c.IsClosed {
		c.Err = errors.New("ToArray() called on closed cursor")
		return result
	}

	err := c.getMongoCursor()

	if err != nil {
		return result
	}

	if len(parm) > 0 {
		err = c.MongoCursor.All(context.Background(), parm[0])
	} else {
		err = c.MongoCursor.All(context.Background(), &result)
	}

	c.MongoCursor = nil
	c.Close()
	c.Err = err

	return result
}

// Count returns a count of the (remaining) documents for the cursor.
func (c *Cursor) Count() int {
	return len(c.ToArray())
}

// Pretty returns a pretty string version of the remaining documents for a cursor.
// Unlike String() it shows the bson.D as key:value instead of {Key:key Value:value}
func (c *Cursor) Pretty() string {
	var buf bytes.Buffer
	docs := c.ToArray()
	for _, doc := range docs {
		buf.WriteString("{ \n")
		for _, v := range doc {
			fmt.Fprintf(&buf, "    %s : %v \n", v.Key, v.Value)
		}
		buf.WriteString("} \n")
	}

	return buf.String()
}

// String fulfills the Stringer interface.
// Calling this will return a string containing the "pretty" print
// contents of the ToArray() function. ToArray() returns an array
// with all of the documents remaining for the cursor and closes the cursor.
func (c *Cursor) String() string {
	json, _ := json.MarshalIndent(c.ToArray(), "", "  ")
	return string(json)
}

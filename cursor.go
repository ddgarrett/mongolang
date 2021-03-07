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
	All others will close the cursor.
	Once the cursor is closed any attempt to use it will cause a panic.

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
		c.Options = options.FindOptions{}
	} else {
		c.Err = errors.New("Close called on already closed cursor")
	}

	return c.Err
}

// bufferNext reads the next document, if it exists, into the Cursor buffer.
// Returns true if there is a next document.
// If no next document, automatically closes the cursor.
// Since this is an internal method it will panic if called with a closed cursor.
func (c *Cursor) bufferNext() bool {
	if c.IsClosed {
		panic("internal error: bufferNext() called after cursor closed")
	}

	if c.MongoCursor == nil {
		cursor, err := c.Collection.MongoColl.Find(context.Background(), c.Filter, &c.Options)
		c.IsClosed = false

		if err != nil {
			c.Close()
			c.Err = err
			return false
		}

		c.MongoCursor = cursor
	}

	hasNext := c.MongoCursor.Next(context.Background())

	if !hasNext {
		c.Close()
		return false
	}

	c.NextDoc = &bson.D{}
	err := c.MongoCursor.Decode(c.NextDoc)

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

// Sort specifies the bson.D to be used to sort the cursor results
func (c *Cursor) Sort(sortSequence bson.D) *Cursor {
	if c.IsClosed {
		c.Err = errors.New("Sort called on closed cursor")
	} else {
		c.Options.Sort = &sortSequence
	}

	return c
}

// Skip specifies the number of documents to skip before returning the first document
func (c *Cursor) Skip(skipCount int64) *Cursor {
	if c.IsClosed {
		c.Err = errors.New("Skip called on closed cursor")
	} else {
		c.Options.Skip = &skipCount
	}

	return c
}

// Limit specifies the max number of documents to return
func (c *Cursor) Limit(limitCount int64) *Cursor {
	if c.IsClosed {
		c.Err = errors.New("Limit called on closed cursor")
	} else {
		c.Err = nil
		c.Options.Limit = &limitCount
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

// ToArray returns all of the remaining documents for a cursor
// in an array.
func (c *Cursor) ToArray() []bson.D {
	result := []bson.D{}

	if c.IsClosed {
		c.Err = errors.New("ToArray() called on closed cursor")
		return result
	}

	if c.MongoCursor == nil {
		cursor, err := c.Collection.MongoColl.Find(context.Background(), c.Filter, &c.Options)

		if err != nil {
			c.Close()
			c.Err = err
			return result
		}

		c.MongoCursor = cursor
		c.IsClosed = false
	}

	err := c.MongoCursor.All(context.Background(), &result)

	c.MongoCursor = nil
	c.Close()
	c.Err = err

	return result
}

// String fulfills the Stringer interface.
// Calling this will return a string containing the "pretty" print
// contents of the ToArray() function. That function returns an array
// with all of the documents remaining for the cursor and closes the cursor.
func (c *Cursor) String() string {
	json, _ := json.MarshalIndent(c.ToArray(), "", "  ")
	return string(json)
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

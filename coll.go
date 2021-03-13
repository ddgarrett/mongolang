package mongolang

/*
	Methods to support access of MongoDB Collections.
*/

import (
	"context"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewCursor creates a new cursor for this collection
func (c *Coll) NewCursor() *Cursor {
	result := Cursor{
		Collection:   c,
		IsClosed:     true,
		IsFindCursor: false,
		FindOptions:  options.FindOptions{},
		AggrOptions:  options.AggregateOptions{},
	}

	return &result
}

// FindOne returns a single MongoDB Document
// All parms are optional. If present, the following parms
// are recognized:
// 	parms[0] - query - bson.M or bson.D defines of which documents to select
//  parms[1] - projection - bson.M or bson.D defines which fields to retrieve
// TODO: process projection parms
func (c *Coll) FindOne(parms ...interface{}) *bson.D {

	//TODO: add processing of project parm

	var filter interface{}

	if len(parms) > 0 {
		filter, c.Err = verifyParm(parms[0], bsonDAllowed|bsonMAllowed)
		c.DB.Err = c.Err
		if c.Err != nil {
			return &bson.D{}
		}
	}

	result := c.MongoColl.FindOne(context.Background(), filter)
	c.DB.Err = result.Err()
	c.Err = result.Err()

	document := bson.D{}
	if result.Err() != nil {
		return &document
	}

	c.DB.Err = result.Decode(&document)
	c.Err = c.DB.Err
	return &document
}

// Find returns a Cursor
// Parms are optional. If present, the following parms
// are recognized:
// 	parms[0] - query - bson.M defines of which documents to select
//  parms[1] - projection - bson.D defines which fields to retrieve
// TODO: process projection parms
func (c *Coll) Find(parms ...interface{}) *Cursor {

	//TODO: add processing of project parm
	//TODO: what to do if incorrect type of parm passed? Maybe return nil?

	result := c.NewCursor()
	result.IsFindCursor = true
	result.IsClosed = false

	if len(parms) > 0 {
		result.Filter, c.Err = verifyParm(parms[0], (bsonDAllowed | bsonMAllowed))
		c.DB.Err = c.Err
		if c.Err != nil {
			// will cause a later error
			// if Find() is chained
			// or c.Err/c.DB.Err is not checked
			return nil
		}
	} else {
		result.Filter = bson.D{}
	}

	return result
}

// Aggregate returns a cursor for an aggregation pipeline operation.
// The pipeline passed can be one of: []bson.D, bson.A, string
// If bson.A, each entry must be a bson.D
// If string, must be a valid JSON doc that parses to a valid bson.A
func (c *Coll) Aggregate(pipeline interface{}, parms ...interface{}) *Cursor {

	//TODO: process other parms
	//TODO: what to do if incorrect type of parm passed? Maybe return nil?

	result := c.NewCursor()
	result.IsFindCursor = false
	result.IsClosed = false

	result.AggrPipeline, c.Err = verifyParm(pipeline, (bsonAAllowed | bsonDSliceAllowed))
	c.DB.Err = c.Err

	return result
}

// given an input aggregation pipeline interface,
// convert it to a []bson.D or bson.A
func getPipeline(in *interface{}) (interface{}, error) {
	switch v := (*in).(type) {
	case []bson.D:
		return v, nil
	case bson.A:
		return v, nil
	case string:
		return parseJSONPipeline(v)
	default:
		// fmt.Println("getPipeline found unrecognized type")
		err := fmt.Errorf("Invalid parm %T, expected []bson.D, bson.A or JSON string", v)
		fmt.Println(err)
		return nil, err
	}
}

// Convert a JSON string to an aggregation pipeline []bson.D
// or bson.A
func parseJSONPipeline(in string) (interface{}, error) {
	parser := JSONToBSON{}
	parser.ParseJSON(in)

	if parser.Err != nil {
		return nil, parser.Err
	}

	if parser.IsBSOND {
		return []bson.D{parser.BSOND}, nil
	}

	return parser.BSONA, nil
}

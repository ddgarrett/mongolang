package mongolang

/*
	Methods to support access of MongoDB Collections.
*/

import (
	"context"
	"errors"

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

	filter := convertBSONParm(0, parms)

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
		filter := convertBSONParm(0, parms)
		if f, ok := filter.(bson.M); ok {
			result.Filter = &f
		} else {
			result.Err = errors.New("Invalid filter. Expected bson.M ")
		}
	} else {
		result.Filter = &bson.M{}
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

	realPipeline, err := getPipeline(&pipeline)
	c.Err = err
	c.DB.Err = err

	result.AggrPipeline = realPipeline

	return result
}

// given an input aggregation pipeline interface,
// convert it to a []bson.D
func getPipeline(in *interface{}) (interface{}, error) {
	switch v := (*in).(type) {
	case []bson.D:
		// fmt.Println("getPipeline found []bson.D")
		return v, nil
	case bson.A:
		// fmt.Println("getPipeline found bson.A")
		return converBSONAPipeline(v)
	case string:
		// fmt.Println("getPipeline found string")
		return parseJSONPipeline(v)
	default:
		// fmt.Println("getPipeline found unrecognized type")
		err := fmt.Errorf("Invalid parm %T, expected []bson.D, bson.A or JSON string", v)
		fmt.Println(err)
		return nil, err
	}
}

// Convert bson.A to []bson.D
// Error if not every entry in bson.A is a bson.D
func converBSONAPipeline(input bson.A) ([]bson.D, error) {
	result := []bson.D{}

	for _, entry := range input {
		switch v := entry.(type) {
		case bson.D:
			result = append(result, entry.(bson.D))
		default:
			// fmt.Println("error in type of bson.A entry in converBSONAPipeline")
			err := fmt.Errorf("Converting bson.A to []bson.D, found type %T", v)
			return result, err
		}
	}
	return result, nil
}

// Convert a JSON string to an aggregation pipeline []bson.D
func parseJSONPipeline(in string) ([]bson.D, error) {
	parser := JSONToBSON{}
	parser.ParseJSON(in)

	if parser.Err != nil {
		return nil, parser.Err
	}

	if parser.IsBSOND {
		return nil, errors.New("Pipeline string not a bson.A (array")
	}

	return converBSONAPipeline(parser.BSONA)
}

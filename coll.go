package mongolang

/*
	Methods to support access of MongoDB Collections.
*/

import (
	"context"
	"errors"

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

// Aggregate returns a cursor for an aggregation pipeline operation
func (c *Coll) Aggregate(parms ...interface{}) *Cursor {
	result := c.NewCursor()

	return result
}

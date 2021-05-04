package mongolang

/*
	Methods to support access of MongoDB Collections.
*/

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// colOkay if Coll is properly linked to a DB
// and the DB is okay.
// NOTE that this does NOT specifically check that there
// haven't been any errors.
// To do that, check that Err() == nil.
func (c *Coll) collOkay() bool {
	if c.DB == nil || !c.DB.dbOkay() {
		return false
	}

	return true
}

// Return any errors or nil if no error
func (c *Coll) Err() error {

	if !c.collOkay() {
		if c.DB == nil {
			return ErrInvalidColl
		}
	}

	return c.DB.Err
}

// If we don't already have an error
// set the error for the related DB.
// This does ensure that we are properly linked
// to a valid DB before trying to set the DB.Err.
func (c *Coll) setErr(err error) {
	if c.Err() == nil {
		c.DB.Err = err
	}
}

// resetErrors resets any errors if the related
// DB is okay.
func (c *Coll) resetErrors() {
	if c.collOkay() {
		c.DB.Err = nil
	}
}

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
func (c *Coll) FindOne(parms ...interface{}) *bson.D {

	if !c.collOkay() {
		return &bson.D{}
	}

	c.resetErrors()

	var filter interface{}
	var err error

	if len(parms) > 0 {
		filter, err = verifyParm(parms[0], bsonDAllowed|bsonMAllowed)
		c.setErr(err)
		if err != nil {
			return &bson.D{}
		}
	} else {
		filter = bson.D{}
	}

	findOneOptions := options.FindOneOptions{}
	if len(parms) > 1 {
		findOneOptions.Projection, err = verifyParm(parms[1], (bsonDAllowed | bsonMAllowed))
		c.setErr(err)
		if err != nil {
			return &bson.D{}
		}
	}

	result := c.MongoColl.FindOne(context.Background(), filter, &findOneOptions)
	c.setErr(result.Err())

	document := bson.D{}
	if result.Err() != nil {
		return &document
	}

	c.DB.Err = result.Decode(&document)
	return &document
}

// Find returns a Cursor
// Parms are optional. If present, the following parms
// are recognized:
// 	parms[0] - query - bson.M defines of which documents to select
//  parms[1] - projection - bson.D defines which fields to retrieve
func (c *Coll) Find(parms ...interface{}) *Cursor {

	var err error

	result := c.NewCursor()
	result.IsFindCursor = true
	result.IsClosed = false

	if !c.collOkay() {
		return result
	}

	c.resetErrors()

	if len(parms) > 0 {
		result.Filter, err = verifyParm(parms[0], (bsonDAllowed | bsonMAllowed))
		c.setErr(err)
		if err != nil {
			result.Filter = bson.D{}
			return result
		}
	} else {
		result.Filter = bson.D{}
	}

	if len(parms) > 1 {
		result.FindOptions.Projection, err = verifyParm(parms[1], (bsonDAllowed | bsonMAllowed))
		c.setErr(err)
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

	var err error

	result := c.NewCursor()
	result.IsFindCursor = false
	result.IsClosed = false

	if !c.collOkay() {
		return result
	}

	c.resetErrors()

	result.AggrPipeline, err = verifyParm(pipeline, (bsonAAllowed | bsonDSliceAllowed))
	c.setErr(err)

	return result
}

// InsertOne inserts one document into the Collection.
// Document must be a bson.D or bson.M.
// TODO: implement insert one options
func (c *Coll) InsertOne(document interface{}, opts ...interface{}) *mongo.InsertOneResult {
	if !c.collOkay() {
		return &mongo.InsertOneResult{}
	}

	c.resetErrors()

	insertDocument, err := verifyParm(document, bsonDAllowed|bsonMAllowed)
	c.DB.Err = err
	if err != nil {
		return &mongo.InsertOneResult{}
	}

	result, insertErr := c.MongoColl.InsertOne(context.Background(), insertDocument)
	c.DB.Err = insertErr

	return result
}

// InsertMany inserts a slice of documents into a Collection.
// Documents must be a slice or bson.A of bson.D documents
// TODO: implement insert one options
func (c *Coll) InsertMany(documents interface{}, opts ...interface{}) *mongo.InsertManyResult {
	if !c.collOkay() {
		return &mongo.InsertManyResult{}
	}

	c.resetErrors()

	insertDocuments, parmErr := verifyParm(documents, interfaceSliceAllowed)
	c.DB.Err = parmErr
	if parmErr != nil {
		return &mongo.InsertManyResult{}
	}

	iDocs := insertDocuments.([]interface{})
	result, insertErr := c.MongoColl.InsertMany(context.Background(), iDocs)
	c.DB.Err = insertErr

	return result
}

// DeleteOne deletes a single document. Note that the filter need not specify a
// single document but only one document will be deleted.
// TODO: implement delete options.
func (c *Coll) DeleteOne(filter interface{}, opts ...interface{}) *mongo.DeleteResult {
	if !c.collOkay() {
		return &mongo.DeleteResult{}
	}

	c.resetErrors()

	deleteFilter, err := verifyParm(filter, bsonDAllowed|bsonMAllowed)
	c.DB.Err = err
	if err != nil {
		return &mongo.DeleteResult{}
	}

	result, deleteErr := c.MongoColl.DeleteOne(context.Background(), deleteFilter)
	c.DB.Err = deleteErr

	return result
}

// DeleteMany can delete many documents with one call as specified by the filter
// TODO: implement delete options.
func (c *Coll) DeleteMany(filter interface{}, opts ...interface{}) *mongo.DeleteResult {
	if !c.collOkay() {
		return &mongo.DeleteResult{}
	}

	c.resetErrors()

	deleteFilter, err := verifyParm(filter, bsonDAllowed|bsonMAllowed)
	c.DB.Err = err
	if err != nil {
		return &mongo.DeleteResult{}
	}

	result, deleteErr := c.MongoColl.DeleteMany(context.Background(), deleteFilter)
	c.DB.Err = deleteErr

	return result
}

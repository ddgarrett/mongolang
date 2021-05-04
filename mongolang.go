/*
	MonGolang - v0.1
	A simple framework to provide simple access to MongoDB.
	Provides capabilities similar to the MongoDB Console.
*/

package mongolang

/*
	Structures and methods for accessing MongoDB Databases and
	initializing the framework.
*/

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB struct defines the fields necessary to track the state of a connection
// to a MongoDB server
type DB struct {
	Client *mongo.Client

	Err error

	Database *mongo.Database
	Name     string
}

var ErrNotConnected = errors.New("not connected to a MongoDB")
var ErrNotConnectedDB = errors.New("not connected to a MongoDB Database")

// Coll represents a collection
type Coll struct {
	DB        *DB
	MongoColl *mongo.Collection
	CollName  string
}

var ErrInvalidColl = errors.New("collection not linked to a properly established db")

// Cursor represents a cursor for a Collection
type Cursor struct {
	Collection   *Coll
	MongoCursor  *mongo.Cursor
	IsClosed     bool
	IsFindCursor bool

	NextDoc *bson.D

	Filter      interface{}
	FindOptions options.FindOptions

	AggrPipeline interface{}
	AggrOptions  options.AggregateOptions
}

var ErrInvalidCursor = errors.New("cursor not linked to a properly established collection")
var ErrClosedCursor = errors.New("call made to closed cursor for a method that requires an open cursor")
var ErrNotFindCursor = errors.New("method call requires a Find() cursor")

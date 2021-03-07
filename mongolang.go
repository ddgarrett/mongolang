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

// Coll represents a collection
type Coll struct {
	DB        *DB
	MongoColl *mongo.Collection
	CollName  string
	Err       error
}

// Cursor represents a cursor for a Collection
type Cursor struct {
	Collection  *Coll
	MongoCursor *mongo.Cursor
	IsClosed    bool

	Err error

	NextDoc *bson.D

	Filter  *bson.M
	Options options.FindOptions
}

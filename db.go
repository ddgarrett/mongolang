package mongolang

/*
	Methods for accessing MongoDB Databases and
	initializing the framework.
*/

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Disconnect disconnects the MongoDB and
// cleans up any other resources, resetting the MonGolang structure
func (mg *DB) Disconnect() {

	if mg.Client != nil {
		if err := mg.Client.Disconnect(context.Background()); err != nil {
			if mg.Err == nil {
				mg.Err = err
			}
		}
	}

	mg.Client = nil
	mg.Database = nil
	mg.Name = ""
}

// clientOkay returns true if the mg.Client is okay
func (mg *DB) clientOkay() bool {
	if mg.Client == nil {
		if mg.Err == nil {
			mg.Err = errors.New("not connected to a MongoDB")
		}
		return false
	}

	return true
}

// checkDBOkay checks if the mg.Client and mg.Database
// are properly initialized
func (mg *DB) dbOkay() bool {
	if !mg.clientOkay() {
		return false
	}

	if mg.Database == nil {
		if mg.Err == nil {
			mg.Err = errors.New("not connected to a MongoDB Database")
		}
		return false
	}
	return true
}

// InitMonGolang initializes the connection
// to the MongoDB Database
func (mg *DB) InitMonGolang(connectionURI string) *DB {
	mg.Disconnect()

	// get MongoDB Client
	mg.Client, mg.Err = mongo.NewClient(options.Client().ApplyURI(connectionURI))

	if mg.Err != nil {
		mg.Client = nil
		return mg
	}

	// Connect to Database
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	mg.Err = mg.Client.Connect(ctx)

	if mg.Err != nil {
		mg.Client = nil
	}

	return mg
}

// Use connects the MongoDB Client to the specified Database.
// The MonGolangDB needs to be inialized via mg.InitMonGolang()
// before calling this method.
func (mg *DB) Use(dbName string) *DB {

	// exit if we don't have an mg.Client
	if !mg.clientOkay() {
		return mg
	}

	mg.Name = dbName
	mg.Database = mg.Client.Database(dbName)
	mg.Err = nil
	return mg
}

// Coll returns a collection for a given name
// If there was a previous error
// don't set coll.MongoColl
func (mg *DB) Coll(collectionName string) *Coll {
	coll := new(Coll)
	coll.DB = mg
	coll.CollName = collectionName

	// return if we don't have a Database or Client
	if !mg.dbOkay() {
		coll.Err = mg.Err
		return coll
	}

	coll.MongoColl = mg.Database.Collection(collectionName, nil)
	return coll
}

// ShowDBs returns a list of Database Names
func (mg *DB) ShowDBs() []string {
	if !mg.clientOkay() {
		var result []string
		return result
	}

	databases, err := mg.Client.ListDatabaseNames(context.Background(), bson.M{})
	mg.Err = err

	return databases
}

// ShowCollections returns a list of collections for current Database
func (mg *DB) ShowCollections() []string {
	if !mg.dbOkay() {
		var result []string
		return result
	}

	collections, err := mg.Database.ListCollectionNames(context.Background(), bson.M{})
	mg.Err = err

	return collections
}

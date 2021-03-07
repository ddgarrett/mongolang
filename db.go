package mongolang

/*
	Methods for accessing MongoDB Databases and
	initializing the framework.
*/

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Disconnect disconnects the MongoDB and
// cleans up any other resources, resetting the MonGolang structure
func (mg *DB) Disconnect() {
	if mg == nil {
		return
	}

	mg.Err = nil

	if mg.Client != nil {
		if err := mg.Client.Disconnect(context.Background()); err != nil {
			fmt.Printf("\n ***** Error disconnecting from MongoDB: %v \n", err)
		}
	}

	mg.Client = nil
}

// InitMonGolang initializes the connection
// to the MongoDB Database
func (mg *DB) InitMonGolang(connectionURI string) *DB {
	if mg == nil {
		mg = new(DB)
	} else {
		mg.Disconnect()
	}

	// get MongoDB Client
	mg.Client, mg.Err = mongo.NewClient(options.Client().ApplyURI(connectionURI))

	if mg.Err != nil {
		return mg
	}

	// Connect to Database
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	mg.Err = mg.Client.Connect(ctx)
	return mg
}

// Use connects the MongoDB Client to the specified Database.
// The MonGolangDB needs to be inialized via mg.InitMonGolang() before calling this method.
func (mg *DB) Use(dbName string) *DB {
	if mg == nil || mg.Client == nil {
		return mg
	}

	mg.Name = dbName
	mg.Database = mg.Client.Database(dbName)
	return mg
}

// Coll returns a collection for a given name
func (mg *DB) Coll(collectionName string) *Coll {
	coll := new(Coll)
	coll.DB = mg
	coll.CollName = collectionName
	coll.MongoColl = mg.Database.Collection(collectionName, nil)
	return coll
}

// ShowDBs returns a list of Database Names
func (mg *DB) ShowDBs() []string {
	databases, err := mg.Client.ListDatabaseNames(context.Background(), bson.M{})
	mg.Err = err

	if err != nil {
		fmt.Printf("ShowDBs error: %v \n", err)
	}

	return databases
}

// ShowCollections returns a list of collections for current Database
func (mg *DB) ShowCollections() []string {
	collections, err := mg.Database.ListCollectionNames(context.Background(), bson.M{})
	mg.Err = err

	if err != nil {
		fmt.Printf("ShowCollections error: %v \n", err)
	}

	return collections
}

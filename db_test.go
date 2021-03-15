package mongolang

import (
	"testing"
)

func TestConnectClose(t *testing.T) {
	var db = DB{}

	// error: not connected to DB
	db.ShowDBs()
	db.Use("quickstart")

	// connect and show DBs and collections
	db.InitMonGolang("mongodb://localhost:27017")

	db.ShowDBs()
	db.Use("quickstart")
	db.ShowCollections()

	db.Coll("zips")

	db.Disconnect()

}

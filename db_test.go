package mongolang

import (
	"testing"
)

func TestConnectAndDisconnect(t *testing.T) {
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

func TestDisconnect(t *testing.T) {
	db := DB{}

	// Connect to db
	db.InitMonGolang("mongodb://localhost:27017").Use("quickstart")

	// Disconnect should reset db struct
	db.Disconnect()

	if db.Client != nil || db.Database != nil || db.Name != "" || db.Err != nil {
		t.Errorf("Disconnect did not reset DB struct:  %+v", db)
	}

	// Disconnect should not override any other errors
	db.InitMonGolang("incorrect://localhost:27017")
	db.Disconnect()

	if db.Err == nil {
		t.Errorf("Disconnect overrode previous error:  %+v", db)
	}
}

func TestInitMonGolangErrors(t *testing.T) {
	db := DB{}

	// Init should detect incorrect URI
	db.InitMonGolang("incorrect://localhost:27017")

	if db.Err == nil || db.Client != nil {
		t.Errorf("InitMonGolang error not handled correctly:  %+v", db)
	}
}
func TestUseError(t *testing.T) {
	db := DB{}

	// Use(dbName) should leave any preceding errors
	db.InitMonGolang("incorrect://localhost:27017").Use("quickstart")

	if db.Err == nil || db.Client != nil {
		t.Errorf("Use() did not hand InitMonGolang() correctly:  %+v", db)
	}

	// just in case there was a connection
	db.Disconnect()

	// Should detect not connected to client
	db = DB{}
	db.Use("quickstart")
	if db.Err == nil {
		t.Errorf("Use() did not generate error:  %+v", db)
	}
}

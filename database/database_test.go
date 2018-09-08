package database

import (
	"testing"

	"github.com/asdine/storm"
)

type testDB struct {
	ID    string `storm:"id"`    // primary key
	Group string `storm:"index"` // this field will be indexed
}

func TestOpenProductionDB(t *testing.T) {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		t.Errorf("Database failed to open: %v", err)
		return
	}
}

func TestAccessDB(t *testing.T) {
	db, err := OpenTestDB("")
	if err != nil {
		t.Errorf("Opening test.db failed with: %v", err)
	}

	dbTest := testDB{
		ID:    "test",
		Group: "testing",
	}

	if err := db.Save(&dbTest); err != nil {
		t.Errorf("Save to db failed: %v", err)
	}
	var testGet testDB
	if err := db.One("ID", "test", &testGet); err != nil {
		t.Errorf("Failed to get: %s", err.Error())
	}
	if err := db.DeleteStruct(&testGet); err != nil {
		t.Errorf("Failed to delete: %s", err.Error())
	}

	db.Close()
	if err := db.Save(&dbTest); err == nil {
		t.Error("Could access db after close")
	}
}

package database

import (
	"testing"

	"github.com/asdine/storm"
)

func TestDB(t *testing.T) {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		t.Errorf("Database failed to open: %v", err)
		return
	}
}

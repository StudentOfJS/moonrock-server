package database

import (
	"github.com/asdine/storm"
)

// OpenTestDB attempts to open access to the test DB and returns a pointer to the db and an error
func OpenTestDB(dbName string) (*storm.DB, error) {
	db, err := storm.Open(dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenDB attempts to open access to the test DB and returns a pointer to the db and an error
func OpenDB() (*storm.DB, error) {
	db, err := storm.Open("my.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

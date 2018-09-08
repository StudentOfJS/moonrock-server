package database

import (
	"github.com/asdine/storm"
)

// OpenTestDB accepts a path "../database/" and returns a pointer to the db or an error
func OpenTestDB(path string) (*storm.DB, error) {
	dbName := path + "test.db"
	db, err := storm.Open(dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenProdDB accepts a path "../database/" and returns a pointer to the db or an error
func OpenProdDB(path string) (*storm.DB, error) {
	dbName := path + "prod.db"
	db, err := storm.Open(dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

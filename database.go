package main

import (
	"fmt"
	"log"

	"github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
)

var (
	// Db is the bolt db connection
	Db *storm.DB
)

// CreateBucket takes a bucket name as string and creates a bucket or error
func CreateBucket(c string) {
	Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(c))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	var err error
	Db, err = storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()
}

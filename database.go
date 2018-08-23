package main

import (
	"fmt"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
)

var (
	// Db is the bolt db connection
	Db  *bolt.DB
	err error
)

func createBuckets(c string) {
	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(c))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	Db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(createBucket("newsletter"))
	// createBucket("users")
	defer db.Close()
}

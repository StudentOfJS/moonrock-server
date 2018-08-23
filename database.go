package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	bolt "github.com/coreos/bbolt"
)

var (
	// Db is the bolt db connection
	Db *bolt.DB
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

// BackupHandleFunc returns backup of db for download
func BackupHandleFunc(w http.ResponseWriter, req *http.Request) {
	err := Db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	var err error
	Db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	CreateBucket("users")

	defer Db.Close()
}

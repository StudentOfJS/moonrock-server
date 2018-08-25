package main

import (
	"log"

	"github.com/asdine/storm"
)

var (
	// Db is the bolt db connection
	Db *storm.DB
)

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	var err error
	Db, err = storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := Db.Init(&Login{}); err != nil {
		log.Fatal(err)
	}
	hash, err := HashPassword(ClientSecret)
	if err != nil {
		log.Fatal(err)
	}
	clientCredentials := Login{
		Password: hash,
		Username: ClientID,
	}
	Db.Save(&clientCredentials)

	defer Db.Close()
}

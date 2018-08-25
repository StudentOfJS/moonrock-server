package main

import (
	"log"

	"github.com/asdine/storm"
)

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}

	hash, err := HashPassword(ClientSecret)
	if err != nil {
		log.Fatal(err)
	}
	loginCredentials := Login{
		Password: hash,
		Username: ClientID,
	}
	clientCredentials := User{
		Group: "client",
		Login: loginCredentials,
	}

	db.Save(&clientCredentials)
	defer db.Close()
}

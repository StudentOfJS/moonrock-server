package main

import (
	"log"

	"github.com/asdine/storm"
)

// HandleDB handles the setup of bolt db
func HandleDB() {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	hash, err := HashPassword(ClientSecret)
	if err != nil {
		log.Fatal(err)
		return
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

	hash, err = HashPassword(TestPass)
	if err != nil {
		log.Fatal(err)
		return
	}
	loginCredentials = Login{
		Password: hash,
		Username: TestUser,
	}
	clientCredentials = User{
		Group: "testing",
		Login: loginCredentials,
	}
	db.Save(&clientCredentials)

}

package db

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
	clientCredentials := User{
		Group:    "client",
		Password: hash,
		Username: ClientID,
	}

	db.Save(&clientCredentials)

	hash, err = HashPassword(TestPass)
	if err != nil {
		log.Fatal(err)
		return
	}

	clientCredentials = User{
		Group:    "testing",
		Password: hash,
		Username: TestUser,
	}
	db.Save(&clientCredentials)

}

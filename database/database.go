package database

import (
	"log"

	"github.com/asdine/storm"
	"github.com/studentofjs/moonrock-server/models"
	"github.com/studentofjs/moonrock-server/secrets"
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

	hash, err := models.HashPassword(secrets.ClientSecret)
	if err != nil {
		log.Fatal(err)
		return
	}
	clientCredentials := models.User{
		Group:    "client",
		Password: hash,
		Username: secrets.ClientID,
	}

	db.Save(&clientCredentials)

	hash, err = models.HashPassword(secrets.TestPass)
	if err != nil {
		log.Fatal(err)
		return
	}

	clientCredentials = models.User{
		Group:    "testing",
		Password: hash,
		Username: secrets.TestUser,
	}
	db.Save(&clientCredentials)

}

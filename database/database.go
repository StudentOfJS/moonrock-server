package database

import (
	"errors"
	"log"

	"github.com/asdine/storm"
	"github.com/studentofjs/moonrock-server/models"
	"github.com/studentofjs/moonrock-server/secrets"
)

var (
	// DB is the production database
	DB *storm.DB
	// TestDB is the test database
	TestDB *storm.DB
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

// AccessDB opens access to the production DB
func AccessDB(o <-chan bool, e chan<- error) {
	var err error
	DB, err = storm.Open("my.db")
	if err != nil {
		e <- errors.New("production database failed to open")
	}
	open := <-o
	if !open {
		DB.Close()
	}
}

// AccessTestDB opens access to the test DB
func AccessTestDB(o <-chan bool, e chan<- error) {
	var err error
	TestDB, err = storm.Open("my.db")
	if err != nil {
		e <- errors.New("Test database failed to open")
	}
	open := <-o
	if !open {
		DB.Close()
	}
}

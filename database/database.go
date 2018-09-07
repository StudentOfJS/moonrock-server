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

// // AccessDB opens access to the production DB
// func AccessDB(done <-chan bool, dbErr chan<- error) {
// 	var err error
// 	DB, err = storm.Open("my.db")
// 	if err != nil {
// 		dbErr <- errors.New("production database failed to open")
// 	}
// 	if <-done {
// 		DB.Close()
// 	}
// }

// OpenTestDB opens access to the test DB
func OpenTestDB() (*storm.DB, error) {
	db, err := storm.Open("test.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

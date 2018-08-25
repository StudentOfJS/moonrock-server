package main

import (
	"log"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

// Subscription stores details for sending emails
type Subscription struct {
	Allowed      bool   `storm:"index"` // this field will be indexed
	Confirmed    bool   // this field will not be indexed`
	Email        string `storm:"unique"`       // this field will be indexed with a unique constraint
	Group        string `storm:"index"`        // this field will be indexed
	NewsLetterID int    `storm:"id,increment"` // primary key with auto increment
	LastNL       int16  `storm:"index"`        // this field will not be indexed
}

// Login struct contains the user login data
type Login struct {
	Password []byte // this field will not be indexed
	Username string `storm:"unique"` // this field will be indexed with a unique constraint
}

// User struct contains all the user data
type User struct {
	Address         string // this field will not be indexed
	CountryCode     string // this field will not be indexed
	EthereumAddress string // this field will not be indexed
	FirstName       string // this field will not be indexed
	Group           string `storm:"index"`        // this field will be indexed
	ID              int    `storm:"id,increment"` // primary key with auto increment
	LastName        string // this field will not be indexed
	Login           `storm:"inline"`
}

// TokenSaleUpdatesHandler - signs up from PUT request with email to newsletter
func TokenSaleUpdatesHandler(c *gin.Context) {
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}
	email := c.PostForm("email")
	if err := EmailValid(email); err != nil {
		c.String(400, "invalid email")
	}
	tokenSaleUpdates := Subscription{
		Allowed:      true,
		Confirmed:    false,
		Email:        email,
		Group:        "token_sale_updates",
		NewsLetterID: 0,
		LastNL:       0,
	}
	if err := db.Save(&tokenSaleUpdates); err == storm.ErrAlreadyExists {
		c.String(400, "already signed up")
	}
	c.String(200, "ok")
	defer db.Close()
}

// LoginHandler accepts a username and a password and returns access token or error
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if err := LoginCheck(username, password); err != nil {
		// @todo: return token to client
		c.String(400, "invalid login")
	} else {
		c.String(200, "ok")
	}
}

// RegisterHandler validates the user signup form and saves to db
func RegisterHandler(c *gin.Context) {
	address := c.PostForm("address")
	country := c.PostForm("country")
	ethereum := c.PostForm("ethereum")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	// Need to consider having a newsletter or not
	// newsletter := c.PostForm("newsletter")
	password := c.PostForm("password")
	username := c.PostForm("username")

	if err := LoginValid(username, password); err != nil {
		c.String(400, "invalid login details")
	}

	if err := UserValid(ethereum, firstname, lastname); err != nil {
		c.String(400, "invalid user details")
	}
	// Generate "hash" to store from username password
	hash, err := HashPassword(password)
	if err != nil {
		// TODO: Properly handle error
		c.String(401, "invalid")
	}

	login := Login{
		Password: hash,
		Username: username,
	}

	user := User{
		Address:         address,
		CountryCode:     country,
		EthereumAddress: ethereum,
		FirstName:       firstname,
		Group:           "public_investor",
		LastName:        lastname,
		Login:           login,
	}
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Save(&user); err == storm.ErrAlreadyExists {
		c.String(400, "already signed up")
	}
	// @todo considering logging in after signup
	c.String(200, "ok")
	defer db.Close()
}

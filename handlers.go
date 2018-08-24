package main

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	email := c.PostForm("email")
	if EmailNotValid(email) {
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
	if err := Db.Save(&tokenSaleUpdates); err == storm.ErrAlreadyExists {
		c.String(400, "already signed up")
	}
	c.String(200, "ok")
}

// LoginHandler accepts a username and a password and returns access token or error
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if LoginCheck(username, password) {
		// @todo: return token to client
		c.String(200, "ok")

	} else {
		c.String(400, "invalid login")
	}
}

// RegisterHandler validates the user signup form and saves to db
func RegisterHandler(c *gin.Context) {
	address := c.PostForm("address")
	country := c.PostForm("country")
	ethereum := c.PostForm("ethereum")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	newsletter := c.PostForm("newsletter")
	password := c.PostForm("password")
	username := c.PostForm("username")

	if LoginNotValid(username, password) {
		c.String(400, "invalid login details")
	}

	if UserNotValid(ethereum, firstname, lastname) {
		c.String(400, "invalid user details")
	}
	// Generate "hash" to store from username password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

	if err := Db.Save(&user); err == storm.ErrAlreadyExists {
		c.String(400, "already signed up")
	}
	// @todo considering logging in after signup
	c.String(200, "ok")
}

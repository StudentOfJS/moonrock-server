package main

import (
	"fmt"

	"github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SubscriptionData stores details for sending emails
type SubscriptionData struct {
	Allowed      bool   `storm:"index"` // this field will be indexed
	Confirmed    bool   // this field will not be indexed`
	Email        string `storm:"unique"`       // this field will be indexed with a unique constraint
	Group        string `storm:"index"`        // this field will be indexed
	NewsLetterID int    `storm:"id,increment"` // primary key with auto increment
	LastNL       int16  `storm:"index"`        // this field will not be indexed
}

// login struct contains the user login data
type Login struct {
	Password []byte // this field will not be indexed
	Username string `storm:"unique"` // this field will be indexed with a unique constraint
}

// User struct contains all the user data
type User struct {
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
	tokenSaleUpdates := SubscriptionData{
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
	if LoginNotValid(username, password) {
		c.String(400, "invalid login")
	}
	var user User
	if err := Db.One("UserName", username, &user); err != nil {
		c.String(400, "invalid login")
	}
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		c.String(401, "invalid login")
	}
	// @todo: return token to client
	c.String(200, "ok")

}

// RegisterUser validates the user signup form and saves to db
func RegisterUser(c *gin.Context) {
	newsletter := c.PostForm("newsletter")
	ethereum := c.PostForm("ethereum")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	password := c.PostForm("password")
	username := c.PostForm("username")

	if LoginNotValid(username, password) {
		c.String(400, "invalid login")
	}

	if UserNotValid(ethereum, firstname, lastname) {
		c.String(400, "invalid user details")
	}

	user := User{
		Email:            email,
		EthereumAddress:  "",
		FirstName:        "",
		Group:            "newsletter",
		ID:               CreateUUID(email),
		LastName:         "",
		SubscriptionData: newsletterData,
		Password:         "",
		Username:         "",
	}
	// Generate "hash" to store from username password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		c.String(401, "invalid")
	}
	Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("login"))
		v := b.Get([]byte(username))
		if v != nil {
			c.String(400, "invalid")
			return fmt.Errorf("username exists: %s", err)
		}
		Db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("login"))
			u := tx.Bucket([]byte("users"))

			err = b.Put([]byte(username), []byte(password))
			if err != nil {
				c.String(400, "error")
				return fmt.Errorf("login creation: %s", err)
			}
			err = u.Put([]byte("first"), []byte("true"))
			if err != nil {
				c.String(200, "error writing KV | n")
				return fmt.Errorf("create kv: %s", err)
			}

			bu, err := b.CreateBucketIfNotExists([]byte(username))
			if err != nil {
				c.String(200, "error  creating username bucket")
				return fmt.Errorf("userBucket: %s", err)
			}
			err = bu.Put([]byte("last"), []byte("0"))
			if err != nil {
				c.String(200, "error writing KV | n")
				return fmt.Errorf("create kv: %s", err)
			}

			return nil
		})
		return nil
	})
	c.String(200, "ok")

}

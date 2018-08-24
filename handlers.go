package main

import (
	"fmt"

	"github.com/satori/go.uuid"

	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// NewsletterData stores details for sending emails
type NewsletterData struct {
	Allowed bool  `storm:"index"` // this field will be indexed
	Welcome bool  // this field will not be indexed`
	LastNL  int16 `storm:"index"` // this field will not be indexed
}

// User struct contains all the user data
type User struct {
	Email           string `storm:"unique"` // this field will be indexed with a unique constraint
	EthereumAddress string // this field will not be indexed
	FirstName       string // this field will not be indexed
	Group           string `storm:"index"` // this field will be indexed
	ID              string `storm:"id"`    // primary key
	LastName        string // this field will not be indexed
	NewsletterData  `storm:"inline"`
	Password        string // this field will not be indexed
	Username        string `storm:"index"`
}

// Newsletter - signs up from PUT request with email to newsletter
func Newsletter(c *gin.Context) {
	email := c.PostForm("email")
	if EmailNotValid(email) {
		c.String(200, "invalid email")
	}
	Db.Update(func(tx *bolt.Tx) error {
		u2, err := uuid.FromString(email)
		if err != nil {
			return fmt.Errorf("uuid went wrong: %s", err)
		}
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("newsletter")))

		if err != nil {
			c.String(200, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}
		u, err := u2.MarshalText()
		if err != nil {
			c.String(200, "error  uuid")
			return fmt.Errorf("uuid: %s", err)
		}
		ub, err := b.CreateBucketIfNotExists([]byte(u))
		if err != nil {
			c.String(200, "error  creating user bucket")
			return fmt.Errorf("userBucket: %s", err)
		}
		err = ub.Put([]byte("email"), []byte(email))
		if err != nil {
			c.String(200, "error writing email")
			return fmt.Errorf("create email: %s", err)
		}
		err = ub.Put([]byte("first"), []byte("true"))
		if err != nil {
			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}
		err = ub.Put([]byte("last"), []byte("0"))
		if err != nil {
			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})

	c.String(200, "ok")
}

// Login accepts a username and a password and returns access token or error
func Login(c *gin.Context) error {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if LoginNotValid(username, password) {
		c.String(400, "invalid login")
		return fmt.Errorf("invalid login")
	}

	// get hashed password from db
	Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("login"))
		hash := b.Get([]byte(username))
		if hash == nil {
			c.String(400, "invalid login")
			return fmt.Errorf("username doesn't exist")
		}
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
			c.String(401, "invalid login")
			return fmt.Errorf("passwords don't match: %s", err)
		}
		// @todo: return token to username
		return nil
	})
	c.String(200, "ok")
	return nil
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

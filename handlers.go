package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/validator.v2"

	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// NewsletterSignup is used for the newsletter signup
type NewsletterSignup struct {
	Email string `json:"email"`
}

// NewsletterData stores details for sending emails
type NewsletterData struct {
	Allowed bool              `json:"allowed"`
	Email   *NewsletterSignup `json:"email"`
	First   bool              `json:"first_email"`
	Last    int16             `json:"last_sent"`
}

// Login is for basic username and password login - may swap for SSO
type LoginDetails struct {
	Username string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `validate:"min=8",max=255`
}

// User is encompases all held user data
type User struct {
	Ethereum       string          `validate"regexp=^0x[a-fA-F0-9]{40}$"`
	FirstName      string          `validate:"min=1",max=255`
	LastName       string          `validate:"min=1",max=255`
	LoginDetails   *LoginDetails   `json:"login_details"`
	NewsletterData *NewsletterData `json:"newsletter-data"`
}

// Newsletter - signs up from PUT request with email to newsletter
func Newsletter(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.String(200, "no email provided")
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
	loginRequest := LoginDetails{Username: username, Password: password}
	if errs := validator.Validate(loginRequest); errs != nil {
		c.String(400, "invalid login")
		return fmt.Errorf("invalid login: %s", errs)
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
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" || !validateEmail(username) {
		c.String(400, "invalid")
	}
	password := c.PostForm("password")
	if password == "" {
		c.String(400, "invalid")
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

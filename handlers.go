package main

import (
	"fmt"

	"github.com/satori/go.uuid"

	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
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
type Login struct {
	Password string `json:"password"`
	Username string `json:"email"`
}

// User is encompases all held user data
type User struct {
	Ethereum       string          `json:"ethaddress"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Login          *Login          `json:"login_details"`
	NewsletterData *NewsletterData `json:"newsletter-data"`
}

func Newsletter(c *gin.Context) {
	if c.PostForm("key") == "" || c.PostForm("value") == "" {
		c.String(200, "no email provided")
	}
	Db.Update(func(tx *bolt.Tx) error {
		u2, err := uuid.FromString(c.PostForm("value"))
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
		ub, err := b.CreateBucketIfNotExists([]byte(c.PostForm("value")))
		if err != nil {
			c.String(200, "error  creating user bucket")
			return fmt.Errorf("userBucket: %s", err)
		}
		err = ub.Put([]byte("id"), []byte(u))
		if err != nil {
			c.String(200, "error writing id")
			return fmt.Errorf("create id: %s", err)
		}
		err = ub.Put([]byte("first"), []byte("true"))
		if err != nil {
			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})

	c.String(200, "ok")
}

func Put(c *gin.Context) {

	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		c.String(200, "no bucket name or key | n")
	}

	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		if err != nil {

			c.String(200, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(c.PostForm("key")), []byte(c.PostForm("value")))

		if err != nil {

			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})

	c.String(200, "ok")

}

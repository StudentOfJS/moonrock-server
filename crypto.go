package main

import (
    "fmt"
    "log"

    "golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	user := c.PostForm("user")
	if user == "" {
		c.String(400, "invalid")
	}
	password := c.PostForm("password")
	if password == "" {
		c.String(400, "invalid")
	}
    // Generate "hash" to store from user password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        // TODO: Properly handle error
        c.String(401, "invalid")
    }
    err = Db.View(func(tx *bolt.Tx) error {
      b := tx.Bucket([]byte("login"))
      v := b.Get([]byte(user))
      if v != nil {
        c.String(400, "invalid")
        return fmt.Errorf("user exists: %s", err)
      }
      return nil
    })

    Db.Update(func(tx *bolt.Tx) error {
      b := tx.Bucket([]byte("login"))
      u := tx.Bucket([]byte("users"))

      err = b.Put([]byte(user), []byte(password))
      if err != nil {
        c.String(400, "error")
        return fmt.Errorf("login creation: %s", err)
      }
      err = u.Put([]byte("first"), []byte("true"))
      if err != nil {
        c.String(200, "error writing KV | n")
        return fmt.Errorf("create kv: %s", err)
      }

      bu, err := b.CreateBucketIfNotExists([]byte(user))
      if err != nil {
        c.String(200, "error  creating user bucket")
        return fmt.Errorf("userBucket: %s", err)
      }
      err = ub.Put([]byte("last"), []byte("0"))
      if err != nil {
        c.String(200, "error writing KV | n")
        return fmt.Errorf("create kv: %s", err)
      }
  
      return nil
    })
  
    c.String(200, "ok")
    fmt.Println("Hash to store:", string(hash))
    // Store this "hash" somewhere, e.g. in your database
    Db.

    // After a while, the user wants to log in and you need to check the password he entered
    userPassword2 := "some user-provided password"
    hashFromDatabase := hash

    // Comparing the password with the hash
    if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(userPassword2)); err != nil {
        // TODO: Properly handle error
        log.Fatal(err)
    }

    fmt.Println("Password was correct!")
}



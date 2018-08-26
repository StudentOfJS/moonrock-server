package main

import (
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

// Subscription stores details for sending emails
type Subscription struct {
	Allowed      bool   `storm:"index"`        // this field will be indexed
	Confirmed    bool   `storm:"index"`        // this field will be indexed
	Email        string `storm:"unique"`       // this field will be indexed with a unique constraint
	Group        string `storm:"index"`        // this field will be indexed
	NewsLetterID int    `storm:"id,increment"` // primary key with auto increment
	LastNL       int16  // this field will not be indexed
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
		c.String(500, "server error")
		return
	}
	email := c.PostForm("email")
	if err := EmailValid(email); err != nil {
		c.String(400, "invalid email")
		return
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
		return
	}
	c.String(200, "ok")
	defer db.Close()
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
		return
	}

	if err := UserValid(ethereum, firstname, lastname); err != nil {
		c.String(400, "invalid user details")
		return
	}
	// Generate "hash" to store from username password
	hash, err := HashPassword(password)
	if err != nil {
		c.String(401, "invalid")
		return
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
		c.String(500, "server error")
		return
	}
	if err := db.Save(&user); err == storm.ErrAlreadyExists {
		c.String(400, "already signed up")
		return
	}

	c.JSON(200, gin.H{
		"status":    "updated",
		"address":   address,
		"country":   country,
		"ethereum":  ethereum,
		"firstName": firstname,
		"lastName":  lastname,
	})

	defer db.Close()
	return
}

// UpdateUserHandler updates user details supplied to API
func UpdateUserHandler(c *gin.Context) {
	address := c.PostForm("address")
	country := c.PostForm("country")
	firstname := c.PostForm("firstname")
	idStr := c.PostForm("id")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		c.String(401, "unauthenticated")
		return
	}
	lastname := c.PostForm("lastname")

	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.String(500, "server error")
		return
	}
	if err := db.Update(&User{
		ID:          id,
		Address:     address,
		CountryCode: country,
		FirstName:   firstname,
		LastName:    lastname,
	}); err != nil {
		c.String(400, "update failed")
		return
	}
	c.JSON(200, gin.H{
		"status":    "updated",
		"address":   address,
		"country":   country,
		"firstName": firstname,
		"lastName":  lastname,
	})
	return
}

// ContributionAddressHandler uses an ID to find user and updates their contribution address
func ContributionAddressHandler(c *gin.Context) {
	ethereum := c.PostForm("ethereum")
	idStr := c.PostForm("id")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		c.String(401, "unauthenticated")
		return
	}
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.String(500, "server error")
		return
	}

	if err := db.UpdateField(&User{ID: id}, "EthereumAddress", ethereum); err != nil {
		c.String(400, "update failed")
		return
	}
	c.JSON(200, gin.H{
		"status":   "updated",
		"ethereum": ethereum,
	})
	return

}

// GetContributionAddress returns the saved address of the user
func GetContributionAddress(c *gin.Context) {
	var user User
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.String(500, "server error")
		return
	}
	idStr := c.PostForm("id")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		c.String(401, "unauthenticated")
		return
	}
	err = db.One("ID", id, &user)
	if err != nil {
		c.String(400, "user doesn't exist")
		return
	}
	c.JSON(200, gin.H{
		"status":   "ok",
		"ethereum": user.EthereumAddress,
	})
}

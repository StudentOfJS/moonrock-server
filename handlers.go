package main

import (
	"log"
	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// Login struct contains the user login data
type Login struct {
	Password []byte // this field will not be indexed
	Username string `storm:"unique"` // this field will be indexed with a unique constraint
}

// Subscription stores details for sending emails
type Subscription struct {
	Allowed      bool   `storm:"index"`        // this field will be indexed
	Confirmed    bool   `storm:"index"`        // this field will be indexed
	Email        string `storm:"unique"`       // this field will be indexed with a unique constraint
	Group        string `storm:"index"`        // this field will be indexed
	NewsLetterID int    `storm:"id,increment"` // primary key with auto increment
	LastNL       int16  // this field will not be indexed
}

// User struct contains all the user data
type User struct {
	Address         string // this field will not be indexed
	Confirmed       bool   // this field will not be indexed
	CountryCode     string // this field will not be indexed
	EthereumAddress string // this field will not be indexed
	FirstName       string // this field will not be indexed
	Group           string `storm:"index"`        // this field will be indexed
	ID              int    `storm:"id,increment"` // primary key with auto increment
	LastName        string // this field will not be indexed
	Login           `storm:"inline"`
	ResetCode       uuid.UUID // this field will not be indexed
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

// ConfirmAccountHandler checks a resetCode against the DB and returns an error string or
func ConfirmAccountHandler(c *gin.Context) {
	resetcode := c.PostForm("resetcode")
	rc, _ := uuid.FromString(resetcode)
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.JSON(500, gin.H{"status": "please try again"})
		return
	}
	var user User
	if err := db.One("ResetCode", rc, &user); err != nil {
		c.JSON(400, gin.H{"status": "invalid token, please signup", "to": "/register"})
	}
	if err := db.UpdateField(&User{ID: user.ID}, "Confirmed", true); err != nil {
		c.JSON(500, gin.H{"status": "please try again"})
	}
	c.JSON(200, gin.H{"status": "account successfully confirmed", "to": "/login"})
}

// ForgotPasswordHandler sends a reset email with unique password reset link
func ForgotPasswordHandler(c *gin.Context) {
	resetcode := uuid.Must(uuid.NewV4())
	rc := resetcode.String()
	username := c.PostForm("username")
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.String(500, "server failure")
		return
	}
	if err := db.UpdateField(&Login{Username: username}, "ResetCode", resetcode); err != nil {
		c.JSON(500, gin.H{"status": "please try again"})
	}

	r := NewRequest([]string{username}, "Moonrock password reset")
	r.Send("templates/reset_template.html", map[string]string{
		"reset":    rc,
		"username": username,
	})
	c.JSON(200, gin.H{"status": "check your email"})
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

// RegisterHandler validates the user signup form and saves to db
func RegisterHandler(c *gin.Context) {
	address := c.PostForm("address")
	country := c.PostForm("country")
	ethereum := c.PostForm("ethereum")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	password := c.PostForm("password")
	resetcode := uuid.Must(uuid.NewV4())
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
		Confirmed:       false,
		CountryCode:     country,
		EthereumAddress: ethereum,
		FirstName:       firstname,
		Group:           "public_investor",
		LastName:        lastname,
		Login:           login,
		ResetCode:       resetcode,
	}
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
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

	r := NewRequest([]string{username}, "Moonrock Account Confirmation")
	r.Send("templates/register_template.html", map[string]string{
		"country":  country,
		"ethereum": ethereum,
		"name":     firstname,
	})
	return
}

// ResetPasswordHandler handles the reset code checking and password change
func ResetPasswordHandler(c *gin.Context) {
	password := c.PostForm("password")
	resetcode := c.PostForm("resetcode")
	username := c.PostForm("username")

	// Generate "hash" from password
	hash, err := HashPassword(password)
	if err != nil {
		c.String(400, "invalid password")
		return
	}

	rc, e := uuid.FromString(resetcode)

	if e != nil {
		c.JSON(400, gin.H{"status": "token expired"})
	}

	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		c.String(500, "server failure")
		return
	}

	var user User
	err = db.One("Username", username, &user)

	if uuid.Equal(user.ResetCode, rc) {
		if err := db.UpdateField(&Login{Username: username}, "Password", hash); err != nil {
			c.JSON(400, gin.H{"status": "invalid token, please reset and try again", "to": "/reset"})
		}
		newResetCode := uuid.Must(uuid.NewV4())
		db.UpdateField(&User{ResetCode: rc}, "ResetCode", newResetCode)
		c.JSON(200, gin.H{"status": "success"})
	}
}

// sendTenWelcomeMails gets up to 10 new subscriptions and sends them each a welcome email
func sendTenWelcomeMails(done chan bool) {
	var receivers []Subscription
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Panic(err)
		return
	}
	defer db.Close()

	err = db.Find("Confirmed", false, &receivers, storm.Limit(10))
	if err != nil {
		log.Panic(err)
	}
	subject := "Moonrock ICO Confirmation"

	for _, receiver := range receivers {
		r := NewRequest([]string{receiver.Email}, subject)
		r.Send("templates/welcome_template.html", map[string]string{"username": "Welcome"})
		err := db.UpdateField(&Subscription{NewsLetterID: receiver.NewsLetterID}, "Confirmation", true)
		if err != nil {
			log.Panic(err)
		}
	}
	done <- true
}

// SendWelcomeEmails checks for new subscriptions twice a day and sends
func SendWelcomeEmails() {
	done := make(chan bool, 1)
	go sendTenWelcomeMails(done)
	<-done
	timer := time.NewTimer(12 * time.Hour)
	<-timer.C
	SendWelcomeEmails()
}

// TokenSaleUpdatesHandler - signs up from PUT request with email to newsletter
func TokenSaleUpdatesHandler(c *gin.Context) {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()

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
		c.String(200, "already signed up")
		return
	}
	c.String(200, "ok")
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

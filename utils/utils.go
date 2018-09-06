package main

import (
	"errors"
	"log"

	"github.com/asdine/storm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

// EmailTest contains validation for an email address
type EmailTest struct {
	Email string `validate:"regexp=^[0-9a-zA-Z]+@[0-9a-zA-Z]+(\\.[0-9a-zA-Z]+)+$"`
}

// LoginTest contains validation for login details
type LoginTest struct {
	Username string `validate:"min=4, max=255, regexp=^[0-9a-zA-Z]+@[0-9a-zA-Z]+(\\.[0-9a-zA-Z]+)+$"`
	Password string `validate:"min=8, max=255"`
}

// UserTest contains validation for user details
type UserTest struct {
	Ethereum  string `validate:"regexp=^0x[a-fA-F0-9]{40}$"`
	FirstName string `validate:"min=1, max=255, regexp=^[a-zA-Z]+$"`
	LastName  string `validate:"min=1, max=255, regexp=^[a-zA-Z]+$"`
}

// LoginValid returns true if validation fails for username or password
func LoginValid(u string, p string) error {
	loginRequest := LoginTest{Username: u, Password: p}
	if errs := validator.Validate(loginRequest); errs != nil {
		return errors.New("invalid login")
	}
	return nil
}

// UserValid returns true if validation fails for user details
func UserValid(e string, f string, l string) error {
	signupRequest := UserTest{Ethereum: e, FirstName: f, LastName: f}
	if errs := validator.Validate(signupRequest); errs != nil {
		return errors.New("invalid user")
	}
	return nil
}

// EmailValid returns true if validation fails for email
func EmailValid(email string) error {
	emailTest := EmailTest{Email: email}
	if errs := validator.Validate(emailTest); errs != nil {
		return errors.New("invalid email")
	}
	return nil
}

// CreateUUID takes a string represemtation of a uuid and returns an uuid and error
func CreateUUID(stringID string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(stringID)
	return id, err
}

// LoginCheck accepts a username and a password and returns true if checks pass
func LoginCheck(u string, p string) error {
	if err := LoginValid(u, p); err != nil {
		return errors.New("invalid login")
	}
	var user User
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		log.Println("error opening DB")
	}
	if err := db.One("Username", u, &user); err != nil {
		return errors.New("invalid login")
	}
	if !user.Confirmed {
		return errors.New("confirm email")
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(p)); err != nil {
		return errors.New("invalid login")
	}

	return nil
}

// HashPassword takes a string and returns a hash or an error
func HashPassword(p string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("hash failed")
	}
	return hash, nil
}

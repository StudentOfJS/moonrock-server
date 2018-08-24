package main

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

// EmailTest contains validation for an email address
type EmailTest struct {
	Email string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
}

// LoginTest contains validation for login details
type LoginTest struct {
	Username string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `validate:"min=8",max=255`
}

// UserTest contains validation for user details
type UserTest struct {
	Ethereum  string `validate"regexp=^0x[a-fA-F0-9]{40}$"`
	FirstName string `validate:"min=1",max=255`
	LastName  string `validate:"min=1",max=255`
}

// LoginNotValid returns true if validation fails for username or password
func LoginNotValid(u string, p string) bool {
	loginRequest := LoginTest{Username: u, Password: p}
	if errs := validator.Validate(loginRequest); errs != nil {
		return true
	}
	return false
}

// UserNotValid returns true if validation fails for user details
func UserNotValid(e string, f string, l string) bool {
	signupRequest := UserTest{Ethereum: e, FirstName: f, LastName: f}
	if errs := validator.Validate(signupRequest); errs != nil {
		return true
	}
	return false
}

// EmailNotValid returns true if validation fails for email
func EmailNotValid(email string) bool {
	emailTest := EmailTest{Email: email}
	if errs := validator.Validate(emailTest); errs != nil {
		return true
	}
	return false
}

// CreateUUID takes an email and return s an id or error
func CreateUUID(email string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(email)
	return id, err
}

// LoginCheck accepts a username and a password and returns true if checks pass
func LoginCheck(u string, p string) bool {
	if LoginNotValid(u, p) {
		return false
	}
	var user User
	if err := Db.One("UserName", u, &user); err != nil {
		return false
	}
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(p)); err != nil {
		return false
	}
	return true
}

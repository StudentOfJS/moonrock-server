package main

import (
	"gopkg.in/validator.v2"
)

// EmailTest contains validation for an email address
type EmailTest struct {
	Email string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
}

// LoginDetails contains validation for login details
type LoginDetails struct {
	Username string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `validate:"min=8",max=255`
}

// UserDetails contains validation for user details
type UserDetails struct {
	Ethereum  string `validate"regexp=^0x[a-fA-F0-9]{40}$"`
	FirstName string `validate:"min=1",max=255`
	LastName  string `validate:"min=1",max=255`
}

// LoginNotValid returns true if validation fails for username or password
func LoginNotValid(username string, password string) bool {
	loginRequest := LoginDetails{Username: username, Password: password}
	if errs := validator.Validate(loginRequest); errs != nil {
		return true
	}
	return false
}

// UserNotValid returns true if validation fails for user details
func UserNotValid(e string, f string, l string) bool {
	signupRequest := UserDetails{Ethereum: e, FirstName: f, LastName: f}
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

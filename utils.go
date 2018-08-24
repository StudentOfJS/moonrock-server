package main

import (
	"gopkg.in/validator.v2"
)

// LoginNotValid returns true if validation fails for username or password
func LoginNotValid(username string, password string) bool {
	loginRequest := LoginDetails{Username: username, Password: password}
	if errs := validator.Validate(loginRequest); errs != nil {
		return true
	}
	return false
}

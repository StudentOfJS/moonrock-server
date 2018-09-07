package models

import (
	"testing"
)

type testLogin struct {
	username string
	password string
	valid    bool
}

var testLogins = []testLogin{
	{username: "test@test.com", password: "12uhuh4rf89J", valid: true},
	{username: "test2@test.com.au", password: "21268634238432dss", valid: true},
	{username: "test@test.co", password: "uhuhsdfs4rf89J", valid: true},
	{username: "test@test.net.au", password: "needtobreakFree", valid: true},
	{username: "@test.com", password: "12uhuh4rf89J", valid: false},
	{username: "test2@testcom", password: "21268634238432dss", valid: false},
	{username: "test@.co", password: "uhuhsdfs4rf89J", valid: false},
	{username: "test@test.net.au", password: "n", valid: false},
}

func TestLoginValid(t *testing.T) {
	for _, login := range testLogins {
		if login.valid {
			if err := LoginValid(login.username, login.password); err != nil {
				t.Fail()
			}
		}
	}
}

func TestLoginInvalid(t *testing.T) {
	for _, login := range testLogins {
		if !login.valid {
			if err := LoginValid(login.username, login.password); err == nil {
				t.Fail()
			}
		}
	}
}

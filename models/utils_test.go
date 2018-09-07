package models

import (
	"testing"
)

type testLogin struct {
	username string
	password string
	valid    bool
}

type testUser struct {
	ethereum  string
	firstname string
	lastname  string
	valid     bool
}

type testEmail struct {
	email string
	valid bool
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

var testUsers = []testUser{
	{ethereum: "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f", firstname: "coin", lastname: "lancer", valid: true},
	{ethereum: "0x595832F8FC6BF59c85C527fEC3740A1b7a361269", firstname: "Power", lastname: "Ledger", valid: true},
	{ethereum: "0x6a068E0287e55149a2a8396cbC99578f9Ad16A31", firstname: "dave", lastname: "saville", valid: true},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "Ben", lastname: "Georzel", valid: true},
	{ethereum: "08511d6c42Bd247D82746c17a3EE", firstname: "terrence", lastname: "phillip", valid: false},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "", lastname: "Morty", valid: false},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "Rick", lastname: "", valid: false},
}

var testEmails = []testEmail{
	{email: "test@test.com", valid: true},
	{email: "test2@test.com.au", valid: true},
	{email: "test@test.co", valid: true},
	{email: "test@test.net.au", valid: true},
	{email: "@test.com", valid: false},
	{email: "test2@testcom", valid: false},
	{email: "test@.co", valid: false},
	{email: "test.com", valid: false},
}

func TestLoginValid(t *testing.T) {
	for _, login := range testLogins {
		if login.valid {
			if err := LoginValid(login.username, login.password); err != nil {
				t.Fail()
			}
		} else {
			if err := LoginValid(login.username, login.password); err == nil {
				t.Fail()
			}
		}
	}
}

func TestUserValid(t *testing.T) {
	for _, user := range testUsers {
		if user.valid {
			if err := UserValid(user.ethereum, user.firstname, user.lastname); err != nil {
				t.Fail()
			}
		} else {
			if err := UserValid(user.ethereum, user.firstname, user.lastname); err == nil {
				t.Fail()
			}
		}
	}
}

func TestEmailValid(t *testing.T) {
	for _, email := range testEmails {
		if email.valid {
			if err := EmailValid(email.email); err != nil {
				t.Fail()
			}
		} else {
			if err := EmailValid(email.email); err == nil {
				t.Fail()
			}
		}
	}
}

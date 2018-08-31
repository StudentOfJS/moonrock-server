package main

import (
	"reflect"
	"testing"

	"github.com/asdine/storm"
)

var e = "test@test.com.au"

/* ------------------- utils ------------------------- */
func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(Password)
	if err != nil {
		t.Errorf("Expected a hash but recieved an error: %s", err)
	} else {
		hashType := reflect.TypeOf(hash).String()
		if hashType != "[]uint8" {
			t.Error("Expected a []uint8, but recieved: " + hashType)
		}
	}
}

func TestLoginValid(t *testing.T) {
	if err := LoginValid(Username, Password); err != nil {
		t.Errorf("Provided valid username and password got: %s", err)
	}
	if err := LoginValid(Username, "x"); err == nil {
		t.Error("Provided invalid password, expected login check to fail, but it passed")
	}
	if err := LoginValid("username_not_valid", Password); err == nil {
		t.Error("Provided invalid username, expected login check to fail but it passed")
	}
}

func TestUserValid(t *testing.T) {
	e, f, l := "0xCaE9eFE97895EF43e72791a10254d6abDdb17Ae9", "Rod", "Lewis"
	if err := UserValid(e, f, l); err != nil {
		t.Errorf("Provided valid user details, but recieved: %s", err)
	}
	if err := UserValid("not_valid", f, l); err == nil {
		t.Error("Provided invalid eth address, but check passed")
	}
	if err := UserValid(e, "12132", l); err == nil {
		t.Error("Provided invalid name, but check passed")
	}
}

func TestEmailValid(t *testing.T) {
	if err := EmailValid(e); err != nil {
		t.Errorf("Provided valid email, but recieved: %s", err)
	}
}

func TestCreateUUID(t *testing.T) {
	s := "F0001234-0451-4000-B000-000000000000"
	id, err := CreateUUID(s)
	if err != nil {
		t.Errorf("Provided valid string, expected id, but recieved: %s", err)
	} else {
		if reflect.TypeOf(id).String() != "uuid.UUID" {
			t.Error("Expected type of id to be uuid.UUID")
		}
	}
	if _, err := CreateUUID("not_valid_string"); err == nil {
		t.Error("Provided invalid string, expected error, but recieved nil")
	}
}

func TestLoginCheck(t *testing.T) {
	if err := LoginCheck(e, "invalid_login"); err == nil {
		t.Error("Provided invalid login details, expected error, but recieved nil")
	}
	if err := LoginCheck(TestUser, TestPass); err.Error() != "confirm email" {
		t.Errorf("Provided valid login details, but recieved: %s", err)
	}
}

/* ------------------- Database ------------------------- */

func TestDB(t *testing.T) {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		t.Errorf("Database failed to open: %v", err)
		return
	}
}

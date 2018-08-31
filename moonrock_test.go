package main

import (
	"reflect"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(Password)
	if err != nil {
		t.Errorf("Expected a hash but recieved an error: %d", err)
	} else {
		hashType := reflect.TypeOf(hash).String()
		if hashType != "[]uint8" {
			t.Error("Expected a []uint8, but recieved: " + hashType)
		}
	}
}

func TestLoginValid(t *testing.T) {
	if err := LoginValid(Username, Password); err != nil {
		t.Errorf("Provided valid username and password got: %d", err)
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
		t.Errorf("Provided valid user details, but recieved: %d", err)
	}
	if err := UserValid("not_valid", f, l); err == nil {
		t.Error("Provided invalid eth address, but check passed")
	}
	if err := UserValid(e, "12132", l); err == nil {
		t.Error("Provided invalid name, but check passed")
	}

}

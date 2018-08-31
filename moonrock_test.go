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

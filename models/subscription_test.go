package models

import "testing"

var signups = []string{
	"test1@test.com",
	"test2@test.com",
	"test3@test.com",
	"test4@test.com",
}

// @todo finish test
func TestTGENewsletter(t *testing.T) {
	for _, email := range signups {
		TGENewsletter(email)
	}
}

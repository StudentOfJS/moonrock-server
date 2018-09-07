package models

import (
	"testing"
)

type testResponses struct {
	responses  []string
	serverCode int
}

var tests = []testResponses{
	{[]string{"updated", "check your email", "ok", "account successfully confirmed"}, 200},
	{[]string{
		"invalid login",
		"invalid email",
		"invalid signup",
		"user doesn't exist",
		"token expired, please try again",
		"update failed",
		"already signed up",
	}, 400},
	{[]string{"unauthenticated"}, 401},
	{[]string{"server error"}, 500},
}

func TestGetResponse(t *testing.T) {
	for _, pair := range tests {
		for _, test := range pair.responses {
			code := getResponse(test)
			if code.serverCode != pair.serverCode {
				t.Fail()
			}
		}
	}
}

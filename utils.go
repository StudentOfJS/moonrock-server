package main

type NewUserRequest struct {
	Username string `validate:"min=5,max=255,regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `validate:"min=8",max=255`
}

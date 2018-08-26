package main

// Config is requried for the email server
type Config struct {
	Email    string
	Password string
	Port     int
	Server   string
}

var config = Config{
	Email:    EmailUser,
	Password: EmailPassword,
	Port:     EmailPort,
	Server:   EmailServer,
}

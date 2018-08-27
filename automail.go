package main

import (
	"log"
	"time"

	"github.com/asdine/storm"
)

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

// sendTenWelcomeMails gets up to 10 new subscriptions and sends them each a welcome email
func sendTenWelcomeMails(done chan bool) {
	var receivers []Subscription
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Panic(err)
		return
	}
	defer db.Close()

	err = db.Find("Confirmed", false, &receivers, storm.Limit(10))
	if err != nil {
		log.Panic(err)
	}
	subject := "Moonrock ICO Confirmation"

	for _, receiver := range receivers {
		r := NewRequest([]string{receiver.Email}, subject)
		r.Send("templates/template.html", map[string]string{"username": "Welcome"})
		err := db.UpdateField(&Subscription{NewsLetterID: receiver.NewsLetterID}, "Confirmation", true)
		if err != nil {
			log.Panic(err)
		}
	}
	done <- true
}

// SendWelcomeEmails checks for new subscriptions twice a day and sends
func SendWelcomeEmails() {
	done := make(chan bool, 1)
	go sendTenWelcomeMails(done)
	<-done
	timer := time.NewTimer(12 * time.Hour)
	<-timer.C
	SendWelcomeEmails()
}

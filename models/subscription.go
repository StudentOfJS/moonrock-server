package models

import (
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/studentofjs/moonrock-server/database"
	"github.com/studentofjs/moonrock-server/mailer"
)

// Subscription stores details for sending emails
type Subscription struct {
	Allowed      bool   `storm:"index"`        // this field will be indexed
	Confirmed    bool   `storm:"index"`        // this field will be indexed
	Email        string `storm:"unique"`       // this field will be indexed with a unique constraint
	Group        string `storm:"index"`        // this field will be indexed
	NewsLetterID int    `storm:"id,increment"` // primary key with auto increment
	LastNL       int16  // this field will not be indexed
}

// TGENewsletter - signs user up to newsletter with a provided email
func TGENewsletter(e string) *Response {
	db, err := database.OpenProdDB("./database/")
	if err != nil {
		// return getResponse("server error")
		return &Response{ServerCode: 500, Response: err.Error()}
	}
	defer db.Close()

	if err := EmailValid(e); err != nil {
		return getResponse("invalid email")
	}
	tokenSaleUpdates := Subscription{
		Allowed:      true,
		Confirmed:    false,
		Email:        e,
		Group:        "token_sale_updates",
		NewsLetterID: 0,
		LastNL:       0,
	}
	if err := db.Save(&tokenSaleUpdates); err == storm.ErrAlreadyExists {
		return getResponse("already signed up")
	}
	return getResponse("ok")
}

// sendTenWelcomeMails gets up to 10 new subscriptions and sends them each a welcome email
func sendTenWelcomeMails(done chan bool) {
	var receivers []Subscription
	db, err := database.OpenProdDB("./database/")
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Find("Confirmed", false, &receivers, storm.Limit(10))
	if err != nil {
		log.Panic(err)
	}
	subject := "Moonrock ICO Confirmation"

	for _, receiver := range receivers {
		r := mailer.NewRequest([]string{receiver.Email}, subject)
		r.Send("templates/email/welcome_template.html", map[string]string{"username": "Welcome"})
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

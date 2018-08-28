package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

// Request is for emails
type Request struct {
	body    string
	from    string
	subject string
	to      []string
}

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

const (
	// MIME provides content-type and charset info to the email client
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// NewRequest returns a pointer to a Request
func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

// Send accepts a template and items to insert and sends email
func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}

func (r *Request) sendMail() bool {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", config.Server, config.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.Email, config.Password, config.Server), config.Email, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

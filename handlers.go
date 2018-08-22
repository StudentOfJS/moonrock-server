package main

// NewsletterSignup is used for the newsletter signup
type NewsletterSignup struct {
	Email string `json:"email"`
}

// Newsletter stores details for sending emails
type NewsletterData struct {
	Allowed bool              `json:"allowed"`
	Email   *NewsletterSignup `json:"email"`
	First   bool              `json:"first_email"`
	Last    int16             `json:"last_sent"`
}

// Login is for basic username and password login - may swap for SSO
type Login struct {
	Password string `json:"password"`
	Username string `json:"email"`
}

// User is encompases all held user data
type User struct {
	Ethereum       string          `json:"ethaddress"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Login          *Login          `json:"login_details"`
	NewsletterData *NewsletterData `json:"newsletter-data"`
}

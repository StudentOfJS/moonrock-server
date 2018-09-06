package secrets

import (
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

var (
	//SecretKey is a secret
	SecretKey string
	//ClientID is used for client auth
	ClientID string
	//ClientSecret is used for client auth
	ClientSecret string
	// EmailServer is used for the mail server
	EmailServer string
	// EmailPort is used for the mail server
	EmailPort int
	// Username is used for the mail server
	Username string
	// Password is used for the mail server
	Password string
	// TestPass is a password used for testing purposes
	TestPass string
	// TestUser is a username used for testing purposes
	TestUser string
)

func init() {
	gotenv.Load()
	SecretKey = os.Getenv("SECRET_KEY")
	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	EmailServer = os.Getenv("SMTP_SERVER")
	EmailPort, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	Username = os.Getenv("EMAIL")
	Password = os.Getenv("EMAIL_PASSWORD")
	TestPass = os.Getenv("TEST_PASS")
	TestUser = os.Getenv("TEST_USER")
}

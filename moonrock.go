package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	cors "gopkg.in/gin-contrib/cors.v1"
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

func apiRouter() {
	r := gin.Default()                             // Init Router
	r.Use(gin.Logger())                            // log to Stdout
	r.Use(gin.Recovery())                          // recover from panics with 500
	r.Use(cors.Default())                          // enable Cross-Origin Resource Sharing
	r.PUT("/confirm", ConfirmAccountHandler)       // confirm user account
	r.POST("/register", RegisterHandler)           // register user account
	r.PUT("/reset_password", ResetPasswordHandler) // reset password action
	r.POST("/tgenews", TokenSaleUpdatesHandler)    // signup to token sale news
	RegisterAPI(r)                                 // register router
	log.Fatal(r.Run(":4000"))                      // log server error
}

func main() {
	fmt.Println("Rest API v1.0")
	HandleDB()
	apiRouter()
}

package main

import (
	"fmt"
	"log"
	"os"

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
)

func init() {
	gotenv.Load()
	SecretKey = os.Getenv("SECRET_KEY")
	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
}

func server() {
	r := gin.Default()        // Init Router
	r.Use(gin.Logger())       // log to Stdout
	r.Use(gin.Recovery())     // recover from panics with 500
	r.Use(cors.Default())     // enable Cross-Origin Resource Sharing
	RegisterAPI(r)            // register router
	log.Fatal(r.Run(":4000")) // log server error
}

func main() {
	fmt.Println("Rest API v1.0")
	HandleDB()
	server()
}

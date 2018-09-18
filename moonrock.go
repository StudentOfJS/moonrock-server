package main

import (
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/studentofjs/moonrock-server/database"
	"github.com/studentofjs/moonrock-server/middleware"
	"github.com/studentofjs/moonrock-server/models"
	"github.com/studentofjs/moonrock-server/secrets"
	cors "gopkg.in/gin-contrib/cors.v1"
)

func apiRouter() {
	r := gin.Default()                                          // Init Router
	r.Use(gin.Logger())                                         // log to Stdout
	r.Use(gin.Recovery())                                       // recover from panics with 500
	r.Use(cors.Default())                                       // enable Cross-Origin Resource Sharing
	r.Use(middleware.LimitConnections(10))                      // limit concurrent connections to 10
	r.LoadHTMLGlob("templates/email/*")                         // pre-load email templates
	r.Use(static.Serve("/", static.LocalFile("./views", true))) // serve static site
	r.Use(gzip.Gzip(gzip.DefaultCompression))                   // use gzip with default compression
	RegisterAPI(r)                                              // register router
	log.Fatal(r.Run(":4000"))                                   // log server error
}

func init() {

	db, err := database.OpenProdDB("./database/")
	if err != nil {
		log.Fatalf("database error:%v", err)
		return
	}
	defer db.Close()

	hash, err := models.HashPassword(secrets.ClientSecret)
	if err != nil {
		log.Fatal(err)
		return
	}
	clientCredentials := models.User{
		Group:    "client",
		Password: hash,
		Username: secrets.ClientID,
	}

	db.Save(&clientCredentials)

	hash, err = models.HashPassword(secrets.TestPass)
	if err != nil {
		log.Fatal(err)
		return
	}

	clientCredentials = models.User{
		Group:    "testing",
		Password: hash,
		Username: secrets.TestUser,
	}
	db.Save(&clientCredentials)
}

func main() {
	apiRouter()
}

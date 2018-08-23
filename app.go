package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	cors "gopkg.in/gin-contrib/cors.v1"
)

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

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleRequests() {
	// Init Router
	r := gin.Default()

	// Test Ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// Route Handlers / Endpoints
	r.POST("/newsletter", Newsletter)

	// log server error
	log.Fatal(http.ListenAndServe(":30000", nil))
}

func main() {
	fmt.Println("Rest API v1.0")
	HandleDB()
	handleRequests()
}

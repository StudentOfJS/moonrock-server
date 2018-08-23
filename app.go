package main

import (
	"fmt"
	"log"

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
	r.PUT("/newsletter", Newsletter)

	// log server error
	log.Fatal(r.Run("4000"))

}

func main() {
	fmt.Println("Rest API v1.0")
	HandleDB()
	handleRequests()
}

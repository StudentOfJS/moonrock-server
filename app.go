package main

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var cache redis.Conn

func main() {
	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/newsletter", Newsletter).Methods("POST")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

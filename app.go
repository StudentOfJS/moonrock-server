package main

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var cache redis.Conn

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", Signin)
	r.HandleFunc("/welcome", Welcome)
	r.HandleFunc("/refresh", Refresh)
	r.HandleFunc("/signup", Signup)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}

	cache = conn
}

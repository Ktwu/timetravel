package main

import (
	"log"
	"net/http"
	"time"

	"github.com/temelpa/timetravel/server"
)

func main() {
	ttServer := server.NewTimeTravelServer()

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      ttServer.Router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s", address)
	log.Fatal(srv.ListenAndServe())
}

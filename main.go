package main

import (
	"log"
	"net/http"
	"time"

	"github.com/temelpa/timetravel/server"
	"github.com/temelpa/timetravel/service"
)

func main() {
	service, err := service.NewSQLiteRecordService(
		"rainbow_test", service.SQLiteRecordServiceSettings{
			ResetOnStart: false,
		})
	if err != nil {
		log.Fatalf("Unable to launch backing service; got error %v", err)
	}
	ttServer := server.NewTimeTravelServer(&service)

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

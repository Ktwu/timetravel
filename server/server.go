package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/api"
	"github.com/temelpa/timetravel/service"
)

type TimeTravelServer struct {
	Router *mux.Router
	Api *api.API
}

func NewTimeTravelServer() TimeTravelServer {
	router := mux.NewRouter()
	service := service.NewInMemoryRecordService()
	api := api.NewAPI(&service)

	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Printf("error: %v", err)
		}
	})
	api.CreateRoutes(apiRoute)

	return TimeTravelServer{router, api}
}
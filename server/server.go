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
	Api    *api.API
}

func NewTimeTravelServer(service service.RecordService) TimeTravelServer {
	router := mux.NewRouter()
	api := api.NewAPI(service)
	api.CreateRoutes(router)

	apiRoute := router.PathPrefix("/api/v{apiVersion:1}").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Printf("error: %v", err)
		}
	})

	return TimeTravelServer{router, api}
}

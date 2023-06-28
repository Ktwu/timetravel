package server

import (
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

	// TODO figure out an intelligent health check. Until then, don't
	// even bother exposing an API for it.

	return TimeTravelServer{router, api}
}

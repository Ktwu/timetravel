package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/entity"
	"github.com/temelpa/timetravel/service"
)

type APIv1 struct {
	records service.RecordServiceV1
}

// generates all api routes
func (a *APIv1) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.getRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.postRecords).Methods("POST")
}

func (a *APIv1) Sanitize(r entity.Record) interface{} {
	return r.IntoV1()
}

func (a *APIv1) getRecords(w http.ResponseWriter, r *http.Request) {
	GetRecords(a, a.records, w, r)
}

func (a *APIv1) postRecords(w http.ResponseWriter, r *http.Request) {
	PostRecords(a, a.records, w, r)
}

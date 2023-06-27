package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/entity"
	"github.com/temelpa/timetravel/service"
)

type APIv2 struct {
	records service.RecordServiceV2
}

// generates all api routes
func (a *APIv2) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.getRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.postRecords).Methods("POST")
	routes.Path("/records/{id}/versions").HandlerFunc(a.getVersionedRecords).Methods("GET")
	routes.Path("/records/{id}/versions/{vid}").HandlerFunc(a.getVersionedRecord).Methods("GET")
}

func (a *APIv2) Sanitize(r entity.Record) interface{} {
	return r
}

func (a *APIv2) getRecords(w http.ResponseWriter, r *http.Request) {
	GetRecords(a, (a.records).(service.RecordServiceV1), w, r)
}

func (a *APIv2) postRecords(w http.ResponseWriter, r *http.Request) {
	PostRecords(a, (a.records).(service.RecordServiceV1), w, r)
}

func (a *APIv2) getVersionedRecord(w http.ResponseWriter, r *http.Request) {
	GetVersionedRecord(a, a.records, w, r)
}

func (a *APIv2) getVersionedRecords(w http.ResponseWriter, r *http.Request) {
	GetVersionedRecords(a, a.records, w, r)
}

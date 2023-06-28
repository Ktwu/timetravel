package api

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/entity"
	"github.com/temelpa/timetravel/service"
)

type APIVersion interface {
	CreateRoutes(*mux.Router)
	Sanitize(entity.Record) interface{}
}

type API struct {
	apiLock  sync.RWMutex
	versions map[string]APIVersion
}

func NewAPI(records service.RecordService) *API {
	return &API{
		versions: map[string]APIVersion{
			"v1": &APIv1{records},
			"v2": &APIv2{records},
		},
	}
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	for key, api := range a.versions {
		apiRoute := routes.PathPrefix("/api/" + key).Subrouter()
		api.CreateRoutes(apiRoute)
	}
}

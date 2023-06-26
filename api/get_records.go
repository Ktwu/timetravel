package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /records/{id}
// GetRecord retrieves the record.
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	apiVersion, err := strconv.ParseInt(vars["apiVersion"], 10, 32)
	if err != nil {
		// This really shouldn't be possible
		// TODO can we have something cleaner here?
		err := writeError(w, "invalid API version", http.StatusBadRequest)
		logError(err)
		return
	}

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := a.records.GetRecord(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	record.Sanitize(int(apiVersion))
	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

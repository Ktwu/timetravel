package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/service"
)

// GET /records/{id}/versions/{vid}
// GetVersionedRecord retrieves the record at the particular version
func GetVersionedRecord(a APIVersion, records service.RecordServiceV2, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	versionId := vars["vid"]

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	vidNumber, err := strconv.ParseInt(versionId, 10, 32)
	if err != nil || vidNumber <= 0 {
		err := writeError(w, "invalid version id; version id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	rwlock := records.GetRWLockForAPI()
	rwlock.RLock()
	defer rwlock.RUnlock()
	record, err := records.GetVersionedRecord(
		ctx,
		int(idNumber),
		int(vidNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist, or vid %d does not exist", idNumber, vidNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, a.Sanitize(record), http.StatusOK)
	logError(err)
}

// GET /records/{id}/versions
// GetVersionedRecords retrieves all versions of a given record.
// The first element will be the oldest version, and the last the newest version.
func GetVersionedRecords(a APIVersion, records service.RecordServiceV2, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	rwlock := records.GetRWLock()
	rwlock.RLock()
	defer rwlock.RUnlock()
	versions, err := records.GetAllRecordVersions(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	sanitizedVersions := make([]interface{}, len(versions))
	for i, v := range versions {
		sanitizedVersions[i] = a.Sanitize(v)
	}

	response := map[string]interface{}{
		"versions": sanitizedVersions,
	}
	err = writeJSON(w, response, http.StatusOK)
	logError(err)
}

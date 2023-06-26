package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NOTE: This code is mostly copied from github.com/gorilla/mux
// and its examples for writing unit tests
func TestServerSanity(t *testing.T) {
	ttServer := NewTimeTravelServer()

	testRecordPath := fmt.Sprintf("/api/v3/records/%d", 42)
	req, err := http.NewRequest("GET", testRecordPath, nil)
	if err == nil {
		t.Errorf("Should have failed grabbing nonexistant API")
	}

	testRecordPath = fmt.Sprintf("/api/v1/records/%d", 42)
	req, err = http.NewRequest("GET", testRecordPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ttServer.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Should have failed to read non-existant record, but got %v", rr.Code)
	}

	// TODO test posting records and verifying the versions of those records
}
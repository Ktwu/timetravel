package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/temelpa/timetravel/service"
)

// NOTE: This code is mostly copied from github.com/gorilla/mux
// and its examples for writing unit tests
func TestServerSanity(t *testing.T) {
	sqlService, err := service.NewSQLiteRecordService(
		"testdata",
		service.SQLiteRecordServiceSettings{ResetOnStart: true},
	)
	if err != nil {
		t.Fatalf("Unable to create service for testing, error %e", err)
	}
	defer func() {
		os.RemoveAll("testdata")
	}()

	ttServer := NewTimeTravelServer(&sqlService)

	// Test that using an unsupported version fails
	req := newTestRequest(t, "GET", fmt.Sprintf("/api/v0/records/%d", 42), nil)
	rr := httptest.NewRecorder()
	ttServer.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Should have failed grabbing nonexistant API, got %v", rr.Code)
	}

	// Test that grabbing a nonexistant record fails
	testRecordPath := fmt.Sprintf("/api/v1/records/%d", 42)
	req = newTestRequest(t, "GET", testRecordPath, nil)
	rr = httptest.NewRecorder()
	ttServer.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Should have failed to read non-existant record, but got %v", rr.Code)
	}

	// Helper to format a request with a JSON payload, then serve that request
	testServeHTTPWithJSON := func(method string, path string, jsonData map[string]interface{}, expectedResponseData map[string]interface{}) {
		rr := httptest.NewRecorder()
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			t.Fatal(err)
		}
		req = newTestRequest(t, method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		ttServer.Router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Failed %s request for %v, got error %v", method, path, rr.Code)
		}
		compareResponseBody(t, rr.Result(), expectedResponseData)
	}

	// Helper to format a request with no payload, then serve that request
	testServeHTTP := func(method string, path string, expectedResponseData map[string]interface{}) {
		testServeHTTPWithJSON(method, path, map[string]interface{}{}, expectedResponseData)
	}

	// NOTE: json.Unmarshal turns numbers into floats as per documentation
	// Test adding a new entry
	testServeHTTPWithJSON(
		"POST",
		testRecordPath,
		map[string]interface{}{
			"hello":  "world",
			"stable": "data",
		},
		map[string]interface{}{
			"id": float64(42),
			"data": map[string]interface{}{
				"hello":  "world",
				"stable": "data",
			},
		},
	)

	// Test that fetching the record just added matches what we expect
	testServeHTTP(
		"GET",
		testRecordPath,
		map[string]interface{}{
			"id": float64(42),
			"data": map[string]interface{}{
				"hello":  "world",
				"stable": "data",
			},
		},
	)

	// Test mutating the record we just added by deleting, adding, and preserving values
	testServeHTTPWithJSON(
		"POST",
		testRecordPath,
		map[string]interface{}{
			"hello":   nil,
			"goodbye": "world",
		},
		map[string]interface{}{
			"id": float64(42),
			"data": map[string]interface{}{
				"stable":  "data",
				"goodbye": "world",
			},
		},
	)

	// Once again, check that fetching the record matches what we expect
	testServeHTTP(
		"GET",
		testRecordPath,
		map[string]interface{}{
			"id": float64(42),
			"data": map[string]interface{}{
				"stable":  "data",
				"goodbye": "world",
			},
		},
	)
}

func newTestRequest(t *testing.T, method string, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func compareResponseBody(t *testing.T, response *http.Response, expectedBody map[string]interface{}) {
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]interface{}
	json.Unmarshal(jsonBody, &body)
	if !cmp.Equal(body["id"], expectedBody["id"]) {
		t.Errorf("Expected %v(%T), got %v(%T)", expectedBody, expectedBody["id"], body, body["id"])
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestLTPHandlerSinglePair(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ltp", ltpHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/api/v1/ltp?pairs=BTCUSD", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if len(response.LTP) != 1 {
		t.Errorf("Expected 1 pair, got %d", len(response.LTP))
	}
}

func TestLTPHandlerMultiplePairs(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ltp", ltpHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/api/v1/ltp?pairs=BTCUSD,BTCCHF,BTCEUR", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if len(response.LTP) != 3 {
		t.Errorf("Expected 3 pairs, got %d", len(response.LTP))
	}
}

func TestLTPHandlerDefaultPairs(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ltp", ltpHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/api/v1/ltp", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if len(response.LTP) != 3 {
		t.Errorf("Expected 3 default pairs, got %d", len(response.LTP))
	}
}

func TestLTPHandlerWithCache(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ltp", ltpHandler).Methods("GET")

	testPairs := "BTCUSD"

	// Step 1: First Request (Cache Miss)
	req1, err := http.NewRequest("GET", "/api/v1/ltp?pairs="+testPairs, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)

	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response1 APIResponse
	if err := json.NewDecoder(rr1.Body).Decode(&response1); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if len(response1.LTP) != 1 || response1.LTP[0].Pair != testPairs {
		t.Errorf("Unexpected response: %+v", response1)
	}

	// Step 2: Second Request (Cache Hit)
	time.Sleep(2 * time.Second) // Ensure the cache is still valid
	req2, err := http.NewRequest("GET", "/api/v1/ltp?pairs="+testPairs, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response2 APIResponse
	if err := json.NewDecoder(rr2.Body).Decode(&response2); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	// Ensure cache response is identical
	if len(response2.LTP) != 1 || response2.LTP[0].Pair != testPairs || response1.LTP[0].Amount != response2.LTP[0].Amount {
		t.Errorf("Cache response mismatch: first=%+v, second=%+v", response1, response2)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock db functions for testing
var testDB = map[int]bool{1: true, 2: true}

func init() {
	// Override DB functions with mocks
	userExists = func(userID int) bool {
		return testDB[userID]
	}

	saveLocationToDatabase = func(loc LocationRequest) error {
		return nil
	}
}

func TestLocationHandler_Success_LocationTest(t *testing.T) {
	loc := LocationRequest{Latitude: 12.34, Longitude: 56.78, UserID: 1}
	body, _ := json.Marshal(loc)

	req := httptest.NewRequest(http.MethodPost, "/location", bytes.NewReader(body))
	w := httptest.NewRecorder()

	locationHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestLocationHandler_UserDoesNotExist(t *testing.T) {
	loc := LocationRequest{Latitude: 12.34, Longitude: 56.78, UserID: 99}
	body, _ := json.Marshal(loc)

	req := httptest.NewRequest(http.MethodPost, "/location", bytes.NewReader(body))
	w := httptest.NewRecorder()

	locationHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestLocationHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/location", nil)
	w := httptest.NewRecorder()

	locationHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Wrapper handler for tests using helper functions
func testLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loc LocationRequest
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !UserExistsHelper(loc.UserID) {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	if err := SaveLocationHelper(loc); err != nil {
		http.Error(w, "Failed to save location", http.StatusInternalServerError)
		return
	}

	resp := Response{Message: "Location received and saved successfully"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func TestLocationHandler_Success(t *testing.T) {
	payload := LocationRequest{Latitude: 10.0, Longitude: 20.0, UserID: 1}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/location", bytes.NewReader(data))
	w := httptest.NewRecorder()

	testLocationHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Location received and saved successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

func TestLocationHandler_InvalidUser(t *testing.T) {
	payload := LocationRequest{Latitude: 10.0, Longitude: 20.0, UserID: -1}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/location", bytes.NewReader(data))
	w := httptest.NewRecorder()

	testLocationHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestLocationHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/location", nil)
	w := httptest.NewRecorder()

	testLocationHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}
}

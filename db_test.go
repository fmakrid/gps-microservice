package main

import "testing"

func TestSaveLocationDB(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{"ValidUser", 1, false},
		{"InvalidUser", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveLocationHelper(LocationRequest{UserID: tt.userID})
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserExistsDB(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		expected bool
	}{
		{"ExistingUser", 1, true},
		{"NonExistingUser", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserExistsHelper(tt.userID)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

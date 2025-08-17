package main

import "fmt"

// Test helpers accessible by all test files
func SaveLocationHelper(loc LocationRequest) error {
	if loc.UserID <= 0 {
		return fmt.Errorf("invalid user id")
	}
	return nil
}

func UserExistsHelper(userID int) bool {
	return userID > 0
}

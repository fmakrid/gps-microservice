package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url" // Import the url package
	"os"
)

// Response defines a simple response structure
type Response struct {
	Message string `json:"message"`
}

// LocationRequest defines the expected payload for POST /location
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    int     `json:"userID"`
}

// Database connection variable
var db *pgx.Conn

// locationHandler handles POST /location
func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loc LocationRequest
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	if !userExists(loc.UserID) {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	// Insert data into PostgreSQL database
	err = saveLocationToDatabase(loc)
	if err != nil {
		http.Error(w, "Failed to save location", http.StatusInternalServerError)
		log.Printf("Error saving location: %v", err)
		return
	}

	resp := Response{Message: "Location received and saved successfully"}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// saveLocationToDatabase inserts location data into PostgreSQL
func saveLocationToDatabase(loc LocationRequest) error {
	// SQL query to insert the data
	sql := `INSERT INTO public.locations (latitude, longitude, user_id) VALUES ($1, $2, $3)`

	// Execute the insert query
	_, err := db.Exec(context.Background(), sql, loc.Latitude, loc.Longitude, loc.UserID)
	return err
}

func userExists(userID int) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err := db.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}
	return exists
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get PostgreSQL credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// URL-encode the password to escape special characters
	encodedPassword := url.QueryEscape(dbPassword)

	// Create the PostgreSQL connection string with the encoded password
	fmt.Printf("Connecting with: postgres://%s:%s@%s:%s/%s\n", dbUser, encodedPassword, dbHost, dbPort, dbName)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, encodedPassword, dbHost, dbPort, dbName)

	// Set up PostgreSQL connection
	db, err = pgx.Connect(context.Background(), connStr)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Handle POST request to /location
	http.HandleFunc("/location", locationHandler)

	host := "localhost"
	port := 8081
	log.Printf("Server running on %s:%d\n", host, port)

	// Ensure the database connection is closed when the program ends
	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			// Handle the error if db.Close() fails
			log.Printf("Error closing database connection: %v", err)
		}
	}(db, context.Background())

	// Start the server
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

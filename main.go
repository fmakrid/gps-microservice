package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
)

var version = "undefined"

// Response represents the response structure returned by the server.
type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// LocationRequest represents a request containing a user's location information.
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    int     `json:"userID"`
}

// Config holds environment-specific configuration
type Config struct {
	AppEnv     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerHost string
	ServerPort string
}

var db *pgx.Conn

func loadConfig() Config {
	// Only load .env in non-production environments
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found (skipping)")
	}

	return Config{
		AppEnv:     os.Getenv("APP_ENV"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		ServerHost: os.Getenv("SERVER_HOST"), // Can be set to 0.0.0.0 in prod
		ServerPort: os.Getenv("SERVER_PORT"), // e.g., 8081
	}
}

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

	if !userExists(loc.UserID) {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	err = saveLocationToDatabase(loc)
	if err != nil {
		http.Error(w, "Failed to save location", http.StatusInternalServerError)
		log.Printf("Error saving location: %v", err)
		return
	}

	resp := Response{Message: "Location received and saved successfully"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func saveLocationToDatabase(loc LocationRequest) error {
	sql := `INSERT INTO public.locations (latitude, longitude, user_id) VALUES ($1, $2, $3)`
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
	// Load configuration
	config := loadConfig()

	// Connect to PostgreSQL
	encodedPassword := url.QueryEscape(config.DBPassword)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DBUser, encodedPassword, config.DBHost, config.DBPort, config.DBName)

	log.Printf("Connecting to DB at %s (env: %s)", config.DBHost, config.AppEnv)

	var err error
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			// Handle the error if db.Close() fails
			log.Printf("Error closing database connection: %v", err)
		}
	}(db, context.Background())

	// Set up HTTP handler
	http.HandleFunc("/location", locationHandler)

	addr := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	log.Printf("Server running on %s", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

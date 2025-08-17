package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var version = "undefined"
var db *pgx.Conn

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

// ---------------------
// DB functions (can be overridden in tests)
// ---------------------
var saveLocationToDatabase = func(loc LocationRequest) error {
	sql := `INSERT INTO public.locations (latitude, longitude, user_id) VALUES ($1, $2, $3)`
	_, err := db.Exec(context.Background(), sql, loc.Latitude, loc.Longitude, loc.UserID)
	return err
}

var userExists = func(userID int) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err := db.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		return false
	}
	return exists
}

// ---------------------
// HTTP handler
// ---------------------
func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loc LocationRequest
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !userExists(loc.UserID) {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	if err := saveLocationToDatabase(loc); err != nil {
		http.Error(w, "Failed to save location", http.StatusInternalServerError)
		log.Printf("Error saving location: %v", err)
		return
	}

	resp := Response{Success: true, Message: "Location received and saved successfully"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// ---------------------
// Config & main
// ---------------------
func loadConfig() Config {
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
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}

func main() {
	config := loadConfig()

	encodedPassword := url.QueryEscape(config.DBPassword)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DBUser, encodedPassword, config.DBHost, config.DBPort, config.DBName)

	var err error
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	defer db.Close(context.Background())

	http.HandleFunc("/location", locationHandler)

	addr := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

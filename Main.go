package main

import (
	// core packages
	"fmt"
	"log"
	"net/http"
	"os"

	// external packages
	"github.com/joho/godotenv"

	// local packages
	"SpotsTest/database"
	"SpotsTest/handlers"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Initialize database
	connStr := getConnStr()
	err = database.InitDB(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	// Initialize routes
	http.HandleFunc("/api/spots", handler.GetSpotsHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

	log.Println("Server started on port 8080")
}

func getConnStr() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	var connStr string
	if dbPassword == "" {
		connStr = fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable",
			dbUser, dbHost, dbName)
	} else {
		connStr = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbName)
	}

	return connStr
}

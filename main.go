package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

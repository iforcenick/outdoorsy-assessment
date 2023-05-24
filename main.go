package main

import (
	"log"
	"net/http"

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

	controller := RentalController{db}

	router := CreateRouter(controller)

	// Run server
	log.Println("server is running on port 8080. press ctrl + c to quit.")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDBConnection(t *testing.T) {
	// Load .env file for db access credentials
	godotenv.Load()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize database
	db, err := InitDB()
	assert.Nil(t, err, "Error should be nil")
	if err == nil {
		defer db.Close()
	}
}

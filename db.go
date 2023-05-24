package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type RentalsDB struct {
	*sql.DB
}

func InitDB() (*RentalsDB, error) {

	dbname := os.Getenv("dbname")
	user := os.Getenv("user")
	password := os.Getenv("password")
	host := os.Getenv("host")
	sslmode := os.Getenv("sslmode")

	dsn := "dbname=" + dbname +
		"\nuser=" + user +
		"\npassword=" + password +
		"\nhost=" + host +
		"\nsslmode=" + sslmode

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &RentalsDB{db}, nil
}

package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

var DBCredentials = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	"localhost", 5432, "admin", "postgres", "simple-chat")

func OpenDB() {
	db, err := sql.Open("postgres", DBCredentials)

	if err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	DB = db
}

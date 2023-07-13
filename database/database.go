package database

import (
	"database/sql"
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"log"
)

var DB *sql.DB

var (
	host     string
	port     string
	user     string
	password string
	dbname   string
)

func OpenDB() {
	DBCredentials := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.EnvData.PGHOST, config.EnvData.PGPORT, config.EnvData.PGUSER, config.EnvData.PGPASSWORD, config.EnvData.PGDATABASE)

	db, err := sql.Open("postgres", DBCredentials)

	if err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	DB = db
}

package database

import (
	"database/sql"
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"log"
)

var DB *sql.DB

func OpenDB() {
	fmt.Println(config.EnvData)
	fmt.Println(config.EnvData.PGHOST)

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

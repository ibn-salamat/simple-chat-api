package database

import (
	"database/sql"
	"fmt"
	"ibn-salamat/simple-chat-api/helpers"
	"log"
)

var DB *sql.DB

var (
	host     = helpers.GetEnvValue("PGHOST")
	port     = helpers.GetEnvValue("PGPORT")
	user     = helpers.GetEnvValue("PGUSER")
	password = helpers.GetEnvValue("PGPASSWORD")
	dbname   = helpers.GetEnvValue("PGDATABASE")
)

var DBCredentials = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

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

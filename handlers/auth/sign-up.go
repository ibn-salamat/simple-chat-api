package auth

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	"log"
	"net/http"
)

var (
	id       int
	username string
)

type response map[string]string

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT * from users")

	for rows.Next() {
		err := rows.Scan(&id, &username)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, username)
	}

	body := response{
		"message": "signUp",
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonBody)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := response{
		"message": "signin",
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonBody)
}

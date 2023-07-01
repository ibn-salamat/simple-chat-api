package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

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

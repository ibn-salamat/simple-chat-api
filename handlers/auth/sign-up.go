package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"log"
	"net/http"
	"strings"
)

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response map[string]string

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var newUser NewUser

	err := decoder.Decode(&newUser)

	if err != nil || strings.Trim(newUser.Email, " ") == "" || strings.Trim(newUser.Password, " ") == "" {
		w.WriteHeader(400)

		jsonResp := response{
			"errorMessage": "Email and password are required!",
		}

		jsonBody, _ := json.Marshal(jsonResp)

		w.Write(jsonBody)

		return
	}

	newUser.Email = strings.Trim(newUser.Email, " ")

	row := database.DB.QueryRow("SELECT email FROM users WHERE email = $1", newUser.Email).Scan(&newUser.Email)

	if row != sql.ErrNoRows {
		w.WriteHeader(400)

		jsonResp := response{
			"errorMessage": "Email already exists.",
		}

		jsonBody, _ := json.Marshal(jsonResp)

		w.Write(jsonBody)

		return
	}

	fmt.Println(newUser.Email)

	// if err != nil {
	// 	fmt.Println(2)
	// }

	// fmt.Println(rows.Scan(&email))

	// jsonBody, err := json.Marshal(body)

	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// }

	// w.Write(jsonBody)
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

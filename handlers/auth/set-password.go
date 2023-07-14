package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	type Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data Data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	data.Email = strings.Trim(data.Email, " ")
	data.Password = strings.Trim(data.Password, " ")

	if err != nil || data.Email == "" || data.Password == "" {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Incorrect json. Required fields: email, password.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	err = database.DB.QueryRow(`SELECT email FROM users WHERE email = $1`, data.Email).Scan(&data.Email)

	if err == nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Password is already set",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	if err != nil && err != sql.ErrNoRows {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	var confirmed bool

	err = database.DB.QueryRow(
		`SELECT 
	confirmed
	FROM users_confirmation 
	WHERE 
	email = $1
	AND 
	confirmed = $2
	`, data.Email, true,
	).Scan(&confirmed)

	if err != nil && err != sql.ErrNoRows {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	if err == sql.ErrNoRows {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Email is not confirmed. Please confirm your email",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	data.Password = string(hash)

	_, err = database.DB.Exec(`
		INSERT INTO users (email, password)
		VALUES ($1, $2)
	`, data.Email, data.Password)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"message": fmt.Sprintf("Password successfully has been created for %s", data.Email),
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)

	go func() {
		_, err = database.DB.Exec(`
		DELETE FROM users_confirmation
		WHERE email = $1 
	`, data.Email)

		if err != nil {
			log.Println(err)
		}
	}()
}

package auth

import (
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/helpers"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"net/http"
	"strings"
	"time"
)

type NewUser struct {
	Email string `json:"email"`
}

func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notFound.NotFound(w, r)
		return
	}

	var newUser NewUser
	decoder := json.NewDecoder(r.Body)
	newUser.Email = strings.Trim(newUser.Email, " ")

	err := decoder.Decode(&newUser)
	if err != nil || newUser.Email == "" {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Email is required!",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	err = database.DB.QueryRow("SELECT email FROM users WHERE email = $1", newUser.Email).Scan(&newUser.Email)

	if err == nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Email already exists.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	confirmationCode := helpers.CreateConfirmationCode()

	_, err = database.DB.Exec(`
	INSERT INTO users_confirmation 
	(email, confirmation_code) 
	VALUES ($1, $2)
	`, newUser.Email, confirmationCode)

	if err != nil {
		errorMessage := err.Error()

		// TODO
		if strings.Contains(errorMessage, "pq: duplicate key value violates unique constraint \"users_confirmation_email_key\"") {
			errorMessage = "We have already sent code confirmation. Please check your email"
		}

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": errorMessage,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	subject := "Subject: Simple-chat Confirmation"
	content := fmt.Sprintf("Your confirmation code: %s", confirmationCode)

	err = tools.SendMail(newUser.Email, subject, content)

	if err != nil {
		// clear db
		database.DB.Exec("DELETE FROM users_confirmation WHERE email = $1", newUser.Email)

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"message": fmt.Sprintf("We have sent code to your email: %s", newUser.Email),
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)

	// add scheduler for delete unvonfirmed email
	// delete temporary user data
	go func() {
		time.Sleep(300 * time.Second)

		database.DB.Exec("DELETE FROM users_confirmation WHERE email = $1", newUser.Email)
	}()
}

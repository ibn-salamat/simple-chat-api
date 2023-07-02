package auth

import (
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/helpers"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type NewUser struct {
	Email string `json:"email"`
}

func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser NewUser
	decoder := json.NewDecoder(r.Body)
	newUser.Email = strings.Trim(newUser.Email, " ")

	err := decoder.Decode(&newUser)
	if err != nil || newUser.Email == "" {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Email is required!",
		})

		w.Write(jsonBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.DB.QueryRow("SELECT email FROM users WHERE email = $1", newUser.Email).Scan(&newUser.Email)

	if err == nil {
		jsonBody, _ := json.Marshal(response{
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

		jsonBody, _ := json.Marshal(response{
			"errorMessage": errorMessage,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	message := []byte(fmt.Sprintf("Subject: Simple-chat Confirmation. \r\n\r\n\n Your code: %s", confirmationCode))
	addr := "smtp.gmail.com: 587"
	auth := smtp.PlainAuth(
		"",
		"n.salamatoff@gmail.com",
		helpers.GetEnvValue("GOOGLE_GMAIL_KEY"),
		"smtp.gmail.com",
	)
	from := "admin@simple-chat.com"
	to := []string{newUser.Email}

	err = smtp.SendMail(
		addr,
		auth,
		from,
		to,
		message,
	)

	if err != nil {
		// clear db
		database.DB.Exec("DELETE FROM users_confirmation WHERE email = $1", newUser.Email)

		jsonBody, _ := json.Marshal(response{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	jsonBody, _ := json.Marshal(response{
		"message": fmt.Sprintf("We have sent code to your email: %s", newUser.Email),
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)

	// delete temporary user data
	go func() {
		time.Sleep(300 * time.Second)

		database.DB.Exec("DELETE FROM users_confirmation WHERE email = $1", newUser.Email)
	}()
}

package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/helpers"
	"net/http"
	"net/smtp"
	"strings"
)

type NewUser struct {
	Email string `json:"email"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var newUser NewUser

	err := decoder.Decode(&newUser)

	if err != nil || strings.Trim(newUser.Email, " ") == "" {
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := response{
			"errorMessage": "Email is required!",
		}

		jsonBody, _ := json.Marshal(jsonResp)

		w.Write(jsonBody)

		return
	}

	newUser.Email = strings.Trim(newUser.Email, " ")

	row := database.DB.QueryRow("SELECT email FROM users WHERE email = $1", newUser.Email).Scan(&newUser.Email)

	if row != sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := response{
			"errorMessage": "Email already exists.",
		}

		jsonBody, _ := json.Marshal(jsonResp)

		w.Write(jsonBody)

		return
	}

	confirmationCode := helpers.CreateConfirmationCode()

	_, err = database.DB.Query(`
	INSERT INTO users_confirmation 
	(email, confirmation_code) 
	VALUES ($1, $2)
	`, newUser.Email, confirmationCode)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)

		jsonResp := response{
			"errorMessage": err.Error(),
		}

		jsonBody, _ := json.Marshal(jsonResp)

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
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := response{
			"errorMessage": err.Error(),
		}

		// clear db
		database.DB.QueryRow("DELETE FROM users_confirmation WHERE email = $1", newUser.Email).Scan(&newUser.Email)

		jsonBody, _ := json.Marshal(jsonResp)

		w.Write(jsonBody)
		return
	}

	w.WriteHeader(http.StatusOK)

	jsonResp := response{
		"message": fmt.Sprintf("We have sent code to your email: %s", newUser.Email),
	}

	jsonBody, _ := json.Marshal(jsonResp)

	w.Write(jsonBody)
}

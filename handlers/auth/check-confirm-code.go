package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"net/http"
	"strings"
)

type Data struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func CheckConfirmCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	data.Email = strings.Trim(data.Email, " ")
	data.Code = strings.Trim(data.Code, " ")

	if err != nil || data.Email == "" || data.Code == "" {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Incorrect json. Required fields: email, code.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	var confirmed bool
	var leftTriesCount int
	var confirmationCode string

	err = database.DB.QueryRow(
		`SELECT confirmed, left_tries_count, confirmation_code
	FROM users_confirmation 
	WHERE email = $1 
	`, data.Email,
	).Scan(&confirmed, &leftTriesCount, &confirmationCode)

	if err == sql.ErrNoRows {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Time for confirmation is up. Please, start from the beginning.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	} else if err != nil {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	fmt.Printf("\nReceiced code: '%s'. Correct code: '%s'", data.Code, confirmationCode)
	fmt.Println(leftTriesCount)
}

// test confirmation
// go func() {
// 	time.Sleep(5 * time.Second)

// 	_, err = database.DB.Query("UPDATE users_confirmation confirmed SET confirmed = true WHERE email = $1", newUser.Email)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }()

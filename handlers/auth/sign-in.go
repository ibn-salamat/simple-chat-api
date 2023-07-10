package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/tools"
	"log"
	"net/http"
	"strings"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Incorrect json. Required fields: email, password.",
			})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	var passwordHash string

	err = database.DB.QueryRow(`
	SELECT password
	FROM
	users
	WHERE
	email = $1
	`, data.Email).Scan(&passwordHash)

	if err != nil {
		errorMessage := err.Error()

		if err == sql.ErrNoRows {
			errorMessage = "User with this email is not found."
		}

		jsonBody, _ := json.Marshal(response{
			"errorMessage": errorMessage,
			})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(data.Password))

	if err != nil {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Incorrect password",
			})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}


	accesToken, err := tools.GenerateJWT(tools.ACCESS_TOKEN_TYPE, "5min", data.Email)

	if err != nil {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Something went wrong",
			})

		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)

		return
	}

	fmt.Println(accesToken)

}

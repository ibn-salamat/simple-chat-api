package auth

import (
	"database/sql"
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notFound.NotFound(w, r)
		return
	}

	type Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data Data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	data.Email = strings.ToLower(strings.Trim(data.Email, " "))
	data.Password = strings.Trim(data.Password, " ")

	if err != nil || data.Email == "" || data.Password == "" {
		jsonBody, _ := json.Marshal(types.ResponseMap{
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

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": errorMessage,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(data.Password))

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Incorrect password",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	token, err := tools.GenerateJWT(tools.ACCESS_TOKEN_TYPE, data.Email)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Something went wrong",
		})

		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)

		return
	}

	_, err = database.DB.Exec(`
			UPDATE users SET token = $1 WHERE email = $2
	`, token, data.Email)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Something went wrong",
		})

		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)

		return

	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"token": token,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data Data

	w.Header().Set("Content-Type", "application/json")

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

	// bcrypt.

}

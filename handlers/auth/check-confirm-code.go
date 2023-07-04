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

	if err != nil {
		errorMessage := err.Error()

		if err == sql.ErrNoRows {
			errorMessage = "Time for confirmation is up or you tried all tries. Please, start from the beginning."
		}

		jsonBody, _ := json.Marshal(response{
			"errorMessage": errorMessage,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	if data.Code != confirmationCode {
		if leftTriesCount > 1 {
			go func() {
				_, err = database.DB.Exec(`
				UPDATE users_confirmation left_tries_count SET left_tries_count = $2 WHERE email = $1
			`, data.Email, leftTriesCount-1)

				if err != nil {
					fmt.Println(err)
				}
			}()

			jsonBody, _ := json.Marshal(response{
				"errorMessage":     "Incorrect code confirmation",
				"left_tries_count": leftTriesCount - 1,
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)
			return
		} else if leftTriesCount == 1 {
			go func() {
				_, err = database.DB.Exec(`
				DELETE FROM users_confirmation
				WHERE email = $1 
			`, data.Email)

				if err != nil {
					fmt.Println(err)
				}
			}()

			jsonBody, _ := json.Marshal(response{
				"errorMessage":     "Try with new code.",
				"left_tries_count": leftTriesCount - 1,
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)
			return
		}

	}

	if data.Code != confirmationCode {
		fmt.Printf("\nReceiced code: '%s'. Correct code: '%s'", data.Code, confirmationCode)
		fmt.Println(leftTriesCount)
	}

}

// test confirmation
// go func() {
// 	time.Sleep(5 * time.Second)

// 	_, err = database.DB.Query("UPDATE users_confirmation confirmed SET confirmed = true WHERE email = $1", newUser.Email)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }()

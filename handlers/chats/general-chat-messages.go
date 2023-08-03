package chats

import (
	"database/sql"
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
)

type Message struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	MessageType    string `json:"message_type"`
	MessageContent string `json:"message_content"`
	CreatedAt      string `json:"created_at"`
}

func GeneralChatMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		notFound.NotFound(w, r)
		return
	}

	_, err := tools.CheckAuthorization(r)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		log.Println(err)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonBody)

		return
	}

	lastMessageDate := r.URL.Query().Get("lastMessageDate")

	if lastMessageDate == "" {
		rows, err := database.DB.Query(`
	SELECT
	id,
	email,
	message_type,
	message_content,
	created_at
	FROM general_chat_messages
	ORDER BY created_at DESC
	LIMIT 10
	`)

		if err != nil {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage": err.Error(),
			})

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonBody)
			return
		}

		messages, err := getMessagesFromRows(rows)
		if err != nil {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage": err.Error(),
			})

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonBody)
			return
		}

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"data": messages,
		})

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBody)
		return
	}

	rows, err := database.DB.Query(`
	SELECT
	id,
	email,
	message_type,
	message_content,
	created_at
	FROM general_chat_messages
	WHERE $1::timestamptz > created_at 
	ORDER BY created_at DESC
	LIMIT 10
	`, lastMessageDate)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	messages, err := getMessagesFromRows(rows)
	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"data": messages,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func getMessagesFromRows(rows *sql.Rows) (*[]Message, error) {
	defer rows.Close()

	messages := make([]Message, 0, 10)
	for rows.Next() {
		message := Message{}
		err := rows.Scan(&message.Id, &message.Email, &message.MessageType, &message.MessageContent, &message.CreatedAt)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return &messages, nil
}

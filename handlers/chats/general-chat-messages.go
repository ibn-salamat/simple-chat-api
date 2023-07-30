package chats

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	MessageType    string `json:"message_type"`
	MessageContent string `json:"message_content"`
	CreatedAt      string `json:"created_at"`
}

func GeneralChatMessages(w http.ResponseWriter, r *http.Request) {
	var page uint
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

	pageString := r.URL.Query().Get("page")

	if pageString == "" {
		page = 0
	} else {
		formattedPage, err := strconv.ParseUint(pageString, 10, 32)

		if err != nil {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage": "Page should be number",
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)
			return
		}

		page = uint(formattedPage)
	}

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
	OFFSET 10 * $1
	`, page)

	defer rows.Close()

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBody)
		return
	}

	messages := make([]Message, 0, 10)

	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.Id, &message.Email, &message.MessageType, &message.MessageContent, &message.CreatedAt)

		if err != nil {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage": err.Error(),
			})

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonBody)
			return
		}

		messages = append(messages, message)
	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"data": messages,
		"page": page,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

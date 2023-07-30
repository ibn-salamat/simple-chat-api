package chats

import (
	"encoding/json"
	"fmt"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
)

func GeneralChatMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notFound.NotFound(w, r)
		return
	}

	claims, err := tools.CheckAuthorization(r)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		log.Println(err)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonBody)

		return
	}

	fmt.Println(claims)

	// id
	// email
	// message_type
	// message_content
	// created_at
}

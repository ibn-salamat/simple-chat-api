package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"time"

	"github.com/gorilla/websocket"
)

func writeMessage(email string, messageType string, messageContent string) {
	jsonBody, _ := json.Marshal(types.ResponseMap{
		"type":    messageType,
		"email":   email,
		"message": messageContent,
		"date":    time.Now().Format(time.RFC3339),
	})

	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, jsonBody)
	}
}

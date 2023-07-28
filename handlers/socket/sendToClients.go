package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"time"

	"github.com/gorilla/websocket"
)

func sendToClients(email string, messageType string, content string) {
	jsonBody, _ := json.Marshal(types.ResponseMap{
		"type":    messageType,
		"email":   email,
		"content": content,
		"date":    time.Now().Format(time.RFC3339),
	})

	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, jsonBody)
	}
}

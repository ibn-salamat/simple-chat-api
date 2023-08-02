package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func sendToClients(email string, messageType string, content string) {
	id := uuid.New()
	date := time.Now().UTC().Format(time.RFC3339)
	jsonBody, _ := json.Marshal(types.ResponseMap{
		"id":      id,
		"type":    messageType,
		"email":   email,
		"content": content,
		"date":    date,
	})

	_, err := database.DB.Exec(`
		INSERT INTO general_chat_messages (id, email, message_type, message_content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, email, messageType, content, date)

	if err != nil {
		log.Println(err)
	}

	for conn := range clients {
		mutex := new(sync.Mutex)
		mutex.Lock()
		defer mutex.Unlock()
		conn.WriteMessage(websocket.TextMessage, jsonBody)
	}
}

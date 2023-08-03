package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"time"

	"github.com/gorilla/websocket"
)

func sendOnlineUsers(done *chan bool) {
	ticker := time.NewTicker(time.Second * 3)

	for {
		select {
		case <-*done:
			return
		case <-ticker.C:
			userEmails := make([]string, 0, len(clients))

			for _, email := range clients {
				userEmails = append(userEmails, email)
			}

			jsonBody, _ := json.Marshal(types.ResponseMap{
				"type":    "onlineUsers",
				"content": userEmails,
			})

			for conn := range clients {
				conn.WriteMessage(websocket.TextMessage, jsonBody)
			}
		}
	}
}

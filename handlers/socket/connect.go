package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	claims, err := CheckAuthorization(r)
	tickerDone := make(chan bool)

	defer closeConnection(connection, claims.Email, &tickerDone)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"type":    "errorMessage",
			"content": err.Error(),
		})

		connection.WriteMessage(websocket.TextMessage, jsonBody)
		return
	}

	// say hello to all
	clients[connection] = claims.Email
	sendToClients(claims.Email, "connection", "Connected")

	go sendOnlineUsers(&tickerDone)

	// read
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		if time.Now().Unix() > claims.ExpiresAt {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"type":    "errorMessage",
				"content": "token is expired",
			})

			connection.WriteMessage(websocket.TextMessage, jsonBody)
			return
		}

		sendToClients(claims.Email, "message", string(message))
	}
}

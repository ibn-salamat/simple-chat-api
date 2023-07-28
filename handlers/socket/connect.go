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

	defer closeConnection(connection, claims.Email)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"type":    "error",
			"message": err.Error(),
		})

		connection.WriteMessage(websocket.TextMessage, jsonBody)
		return
	}

	// say hello to all
	clients[connection] = claims.Email
	writeMessage(claims.Email, "connection", "Connected")

	// read
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		if time.Now().Unix() > claims.ExpiresAt {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"type":    "error",
				"message": "token is expired",
			})

			connection.WriteMessage(websocket.TextMessage, jsonBody)
			return
		}

		writeMessage(claims.Email, "message", string(message))
	}
}

func closeConnection(connection *websocket.Conn, email string) {
	connection.Close()
	delete(clients, connection)

	if email == "" {
		return
	}

	writeMessage(email, "connection", "Disconnected")
}

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
	defer connection.Close()

	claims, err := CheckAuthorization(r)

	if err != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": err.Error(),
		})

		connection.WriteMessage(websocket.TextMessage, jsonBody)
		return
	}

	clients[connection] = claims.Email
	defer delete(clients, connection)

	// say hello to all
	jsonBody, _ := json.Marshal(types.ResponseMap{
		"email":   claims.Email,
		"message": "Connected",
		"date":    time.Now().Format(time.RFC3339),
	})

	go writeMessage(jsonBody, connection)

	// read
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		if time.Now().Unix() > claims.ExpiresAt {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage": "token is expired",
			})

			connection.WriteMessage(websocket.TextMessage, jsonBody)
			return
		}

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"email":   claims.Email,
			"message": string(message),
			"date":    time.Now().Format(time.RFC3339),
		})

		go writeMessage(jsonBody, connection)
	}

}

func writeMessage(message []byte, currentConn *websocket.Conn) {
	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

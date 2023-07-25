package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

var clients []websocket.Conn

func SocketHandler(ws *websocket.Conn) {
	currentUserEmail, authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	} else {
		log.Println("User connected")
	}

	clients = append(clients, *ws)

	for index, client := range clients {
		email, authorizationError := CheckAuthorization(&client)

		if authorizationError != nil {
			clients = append(clients[:index], clients[index+1:]...)
			return
		}

		if currentUserEmail != email {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"email":   email,
				"message": "connected",
				"date":    time.Now().Format(time.RFC3339),
			})

			websocket.Message.Send(&client, string(jsonBody))
		}
	}

	defer HandleReceive(ws)
}

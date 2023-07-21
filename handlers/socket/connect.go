package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"log"

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

	for _, client := range clients {
		email, authorizationError := CheckAuthorization(&client)

		if authorizationError != nil {
			// delete from connections
			return
		}

		if currentUserEmail != email {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"email":   email,
				"message": "connected",
			})

			websocket.Message.Send(&client, string(jsonBody))
		}
	}

	defer HandleReceive(ws)
}

package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"golang.org/x/net/websocket"
)

var clients []websocket.Conn

var currentUserEmail string

func SocketHandler(ws *websocket.Conn) {
	email, authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	} else {
		log.Println("User connected")
	}

	currentUserEmail = email

	clients = append(clients, *ws)

	for _, client := range clients {
		email, authorizationError := CheckAuthorization(&client)

		if authorizationError != nil {
			// delete from connections
			return
		}

		if email != currentUserEmail {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"email": email,
				"message": "connected",
			})

			websocket.Message.Send(&client, string(jsonBody))
		}
	}

	defer HandleReceive(ws)
}
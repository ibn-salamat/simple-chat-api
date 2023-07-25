package socket

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func HandleReceive(ws *websocket.Conn) {
	for {
		var reply string
		currentUserEmail, authorizationError := CheckAuthorization(ws)

		if authorizationError != nil {
			for index, client := range clients {
				if client == ws {
					clients = append(clients[:index], clients[index+1:]...)
				}
			}

			log.Println("User disconnected")
			return
		}

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			if err.Error() == "EOF" {
				err = ws.Close()
				if err != nil {
					log.Println(err)
				}

				for index, client := range clients {
					if client == ws {
						clients = append(clients[:index], clients[index+1:]...)
					}
				}

				log.Println("User disconnected")
				break
			}

			log.Println(err)
		}

		for index, client := range clients {
			email, authorizationError := CheckAuthorization(client)

			if authorizationError != nil {
				clients = append(clients[:index], clients[index+1:]...)
				return
			}

			if currentUserEmail != email {
				jsonBody, _ := json.Marshal(types.ResponseMap{
					"email":   email,
					"message": reply,
					"date":    time.Now().Format(time.RFC3339),
				})

				websocket.Message.Send(client, string(jsonBody))
			}
		}
	}
}

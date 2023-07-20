package socket

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"ibn-salamat/simple-chat-api/types"
	"log"
)

func HandleReceive(ws *websocket.Conn) {
	for {
		var reply string
		_, authorizationError := CheckAuthorization(ws)

		if authorizationError != nil {
			log.Println("User disconnected")
			return
		}

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			if err.Error() == "EOF" {
				err = ws.Close();
				if err != nil {
					log.Println(err)
				};

				log.Println("User disconnected");
				break
			};

			log.Println(err)
		};


		for _, client := range clients {
			email, authorizationError := CheckAuthorization(&client)

			if authorizationError != nil {
				// delete from connections
				return
			}

			if email == currentUserEmail {
				break
			}

			jsonBody, _ := json.Marshal(types.ResponseMap{
				"email": email,
				"message": reply,
			})

			websocket.Message.Send(&client, string(jsonBody))
		}
	}
}
package socket

import (
	"golang.org/x/net/websocket"
	"log"
)

func HandleReceive(ws *websocket.Conn) {
	clients = append(clients, *ws)

	for {
		var reply string
		authorizationError := CheckAuthorization(ws)

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
			if &client == ws {
				return
			}
			websocket.Message.Send(&client, reply)
		}
	}
}
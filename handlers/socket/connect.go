package socket

import (
	"encoding/json"
	"fmt"
	"ibn-salamat/simple-chat-api/tools"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type response map[string]string

func SocketHandler(ws *websocket.Conn) {
	token, err := ws.Request().Cookie("token")

	// check token exists
	if err == http.ErrNoCookie {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": "Required token",
		})

		_, err = ws.Write(jsonBody)

		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	// check token is valid
	err = tools.CheckToken(token.Value)

	if err != nil {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": err.Error(),
			})

		_, err = ws.Write(jsonBody)

		return
	}
	
	log.Println("User connected")

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			if err.Error() == "EOF" {
				err = ws.Close();
				if err != nil {
					log.Println(err)
				}

				log.Println("User disconnected");
				break
			};

			log.Println(err)
		};

		msg := "Received:  " + reply;
		fmt.Println("Sending to client: " + msg)
	}

}

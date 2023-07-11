package socket

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func SocketHandler(ws *websocket.Conn) {
	var err error

	log.Println("User connected")

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			if err.Error() == "EOF" {
				err = ws.Close()
				if err != nil {
					log.Println(err)
				}

				log.Println("User disconnected")
				break
			}

			log.Println(err)
		}

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)
	}

}

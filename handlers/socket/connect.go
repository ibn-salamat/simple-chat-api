package socket

import (
	"fmt"

	"golang.org/x/net/websocket"
)

func SocketHandler(ws *websocket.Conn) {
	var err error

	fmt.Println("user")

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
		}

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)
	}

}

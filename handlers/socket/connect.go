package socket

import (
	"log"
	"golang.org/x/net/websocket"
)

func SocketHandler(ws *websocket.Conn) {
	authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	} else {
		log.Println("User connected")
	}

	defer HandleReceive(ws)
}
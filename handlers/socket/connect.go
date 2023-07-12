package socket

import (
	"golang.org/x/net/websocket"
	"log"
)

type response map[string]string
func SocketHandler(ws *websocket.Conn) {
	authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	}

	log.Println("User connected")

	HandleReceive(ws)
}

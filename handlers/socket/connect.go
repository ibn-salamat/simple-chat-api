package socket

import (
	"log"
	"golang.org/x/net/websocket"
)

var clients []websocket.Conn

func SocketHandler(ws *websocket.Conn) {
	authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	} else {
		log.Println("User connected")
	}

	clients = append(clients, *ws)
	defer HandleReceive(ws)
}
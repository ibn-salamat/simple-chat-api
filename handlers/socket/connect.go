package socket

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"time"
)

type response map[string]string

func SocketHandler(ws *websocket.Conn) {
	authorizationError := CheckAuthorization(ws)

	if authorizationError != nil {
		return
	} else {
		log.Println("User connected")
	}

	defer HandleReceive(ws)

	// test send message
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <- ticker.C:
				// TODO
				authorizationError = CheckAuthorization(ws)

				if authorizationError != nil {
					ticker.Stop()
					return
				}

				if err := websocket.Message.Send(ws, fmt.Sprintf("Message from server: %i", rand.Int())); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

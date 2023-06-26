package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func socketHandler(ws *websocket.Conn) {
	var err error

	fmt.Println("user")

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)
	}
}

func main() {
	http.Handle("/", websocket.Handler(socketHandler))

	fmt.Println("Server started")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

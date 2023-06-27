package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const PORT = ":3000"

func socketHandler(ws *websocket.Conn) {
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

func main() {
	http.Handle("/ws", websocket.Handler(socketHandler))

	fmt.Printf("Server started on PORT %s", PORT)

	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

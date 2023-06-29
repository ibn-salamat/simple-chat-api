package main

import (
	"fmt"
	"ibn-salamat/simple-chat-api/handlers/auth"
	"ibn-salamat/simple-chat-api/handlers/socket"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const PORT = ":3000"

func main() {
	http.Handle("/ws", websocket.Handler(socket.SocketHandler))
	http.Handle("/auth/sign-up", http.HandlerFunc(auth.SignUpHandler))
	http.Handle("/auth/sign-in", http.HandlerFunc(auth.SignInHandler))

	fmt.Printf("Server started on PORT %s", PORT)

	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

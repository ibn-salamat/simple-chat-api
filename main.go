package main

import (
	"database/sql"
	"fmt"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/handlers/auth"
	"ibn-salamat/simple-chat-api/handlers/socket"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"golang.org/x/net/websocket"
)

const PORT = ":3000"

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "admin", "postgres", "simple-chat")

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Unable to connect to Database:", err)
	}

	database.DB = db

	defer db.Close()

	http.Handle("/ws", websocket.Handler(socket.SocketHandler))
	http.Handle("/api/auth/sign-up", http.HandlerFunc(auth.SignUpHandler))
	http.Handle("/api/auth/sign-in", http.HandlerFunc(auth.SignInHandler))

	fmt.Printf("Server started on PORT %s \n", PORT)

	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/handlers/auth"
	"ibn-salamat/simple-chat-api/handlers/socket"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"golang.org/x/net/websocket"
)

const PORT = ":3000"

func init() {
	env, err := godotenv.Read()

	if err != nil {
		log.Fatalf("Could not find .env file")
	}

	config.EnvData.ACCESS_TOKEN_SECRET = env["ACCESS_TOKEN_SECRET"]
	config.EnvData.REFRESH_TOKEN_SECRET = env["REFRESH_TOKEN_SECRET"]

	config.EnvData.GOOGLE_GMAIL_KEY = env["GOOGLE_GMAIL_KEY"]

	config.EnvData.PGDATABASE = env["PGDATABASE"]
	config.EnvData.PGHOST = env["PGHOST"]
	config.EnvData.PGPASSWORD = env["PGPASSWORD"]
	config.EnvData.PGPORT = env["PGPORT"]
	config.EnvData.PGUSER = env["PGUSER"]
}

func main() {
	database.OpenDB()
	defer database.DB.Close()

	http.HandleFunc("/ws",
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(socket.SocketHandler)}
			s.ServeHTTP(w, req)
		})

	http.Handle("/api/auth/sign-up/check-email", http.HandlerFunc(auth.CheckEmailHandler))
	http.Handle("/api/auth/sign-up/check-confirm-code", http.HandlerFunc(auth.CheckConfirmCodeHandler))
	http.Handle("/api/auth/sign-up/set-password", http.HandlerFunc(auth.SetPasswordHandler))
	http.Handle("/api/auth/sign-in", http.HandlerFunc(auth.SignInHandler))

	fmt.Printf("Server started on PORT %s \n", PORT)

	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

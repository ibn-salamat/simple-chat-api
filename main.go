package main

import (
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"ibn-salamat/simple-chat-api/database"
	"ibn-salamat/simple-chat-api/handlers/auth"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/handlers/socket"
	"ibn-salamat/simple-chat-api/middlewares"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func init() {
	env, err := godotenv.Read("./production.env")

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

	config.EnvData.PORT = env["PORT"]

	if config.EnvData.PORT == "" {
		config.EnvData.PORT = "3000"
	}
}

func main() {
	database.OpenDB()
	defer database.DB.Close()

	http.Handle("/ws", http.HandlerFunc(socket.SocketHandler))

	http.Handle("/api/auth/sign-up/check-email", middlewares.SetBasicHeaders(http.HandlerFunc(auth.CheckEmailHandler)))
	http.Handle("/api/auth/sign-up/check-confirm-code", middlewares.SetBasicHeaders(http.HandlerFunc(auth.CheckConfirmCodeHandler)))
	http.Handle("/api/auth/sign-up/set-password", middlewares.SetBasicHeaders(http.HandlerFunc(auth.SetPasswordHandler)))
	http.Handle("/api/auth/sign-in", middlewares.SetBasicHeaders(http.HandlerFunc(auth.SignInHandler)))
	http.Handle("/", middlewares.SetBasicHeaders(http.HandlerFunc(notFound.NotFound)))

	fmt.Printf("Server started on PORT %s \n", config.EnvData.PORT)

	if err := http.ListenAndServe("0.0.0.0:"+config.EnvData.PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

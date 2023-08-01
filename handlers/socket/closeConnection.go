package socket

import "github.com/gorilla/websocket"

func closeConnection(connection *websocket.Conn, email string, tickerDone *chan bool) {
	connection.Close()
	*tickerDone <- true
	delete(clients, connection)

	if email == "" {
		return
	}

	sendToClients(email, "connection", "Disconnected")
}

package socket

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"ibn-salamat/simple-chat-api/tools"
	"log"
	"net/http"
)

func CheckAuthorization(ws *websocket.Conn) error {
	token, resultErr := ws.Request().Cookie("token")
	
	// check token exists
	if resultErr != nil {
		if resultErr == http.ErrNoCookie {
			resultErr = errors.New("token is not exist in cookies")
		}
	}
	
	resultErr = tools.CheckToken(token.Value)
	
	if resultErr != nil {
		jsonBody, _ := json.Marshal(response{
			"errorMessage": resultErr.Error(),
		})

		_, err := ws.Write(jsonBody)
		if err != nil {
			log.Println(err)
		}

		err = ws.Close()
		if err != nil {
			log.Println(err)
		}
	}
	
	return resultErr
}
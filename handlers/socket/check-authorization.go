package socket

import (
	"encoding/json"
	"errors"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"

	"golang.org/x/net/websocket"
)

func CheckAuthorization(ws *websocket.Conn) (string, error) {
	token := ws.Request().URL.Query().Get("token")
	var resultErr error

	if token == "" {
		resultErr = errors.New("token is not exist in params")
	}

	// check token exists
	if resultErr != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
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

		return "", resultErr
	}

	claims, resultErr := tools.CheckToken(token)

	if resultErr != nil {
		jsonBody, _ := json.Marshal(types.ResponseMap{
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

	return claims.Email, resultErr
}

package socket

import (
	"encoding/json"
	"errors"
	"ibn-salamat/simple-chat-api/tools"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func CheckAuthorization(ws *websocket.Conn) error {
	token, resultErr := ws.Request().Cookie("token")

	// check token exists
	if resultErr != nil {
		if resultErr == http.ErrNoCookie {
			resultErr = errors.New("token is not exist in cookies")
		}

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

		return resultErr
	}

	resultErr = tools.CheckToken(token.Value)

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

	return resultErr
}

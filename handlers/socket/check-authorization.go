package socket

import (
	"errors"
	"ibn-salamat/simple-chat-api/tools"
	"net/http"
)

func CheckAuthorization(r *http.Request) (tools.Claims, error) {
	token := r.URL.Query().Get("token")

	var err error

	if token == "" {
		err = errors.New("token is not exist in params")

		return tools.Claims{}, err
	}

	claims, err := tools.CheckToken(token)

	return claims, err
}

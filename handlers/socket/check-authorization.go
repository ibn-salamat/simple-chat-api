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
		tokenCookie, err := r.Cookie("token")

		if err != nil {
			err = errors.New("token does not exist neither in params nor in cookies")
			return tools.Claims{}, err
		}

		token = tokenCookie.Value
	}

	claims, err := tools.CheckToken(token)

	return claims, err
}

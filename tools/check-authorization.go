package tools

import (
	"errors"
	"net/http"
)

func CheckAuthorization(r *http.Request) (Claims, error) {
	token := r.URL.Query().Get("token")

	var err error

	if token == "" {
		tokenCookie, err := r.Cookie("token")

		if err != nil {
			err = errors.New("token does not exist neither in params nor in cookies")
			return Claims{}, err
		}

		token = tokenCookie.Value
	}

	claims, err := CheckToken(token)

	return claims, err
}

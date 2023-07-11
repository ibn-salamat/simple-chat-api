package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"ibn-salamat/simple-chat-api/helpers"
	"log"
)

const REFRESH_TOKEN_TYPE = "REFRESH_TOKEN_TYPE"
const ACCESS_TOKEN_TYPE = "ACCESS_TOKEN_TYPE"


func GenerateJWT(tokenType string , exp string, email string) (string, error) {
	var secretKey string
	
	if tokenType == ACCESS_TOKEN_TYPE {
		secretKey = helpers.GetEnvValue("ACCESS_TOKEN_SECRET")
	} else if tokenType == REFRESH_TOKEN_TYPE {
		secretKey = helpers.GetEnvValue("REFRESH_TOKEN_SECRET")
	} else {
		log.Println("Wrong token type")
		return "", errors.New("Wrong token type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		exp: exp,
		email: email,
	})

	
	tokenString, err := token.SignedString([]byte(secretKey))
	
	return tokenString, err
}
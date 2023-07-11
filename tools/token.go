package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"ibn-salamat/simple-chat-api/helpers"
	"log"
	"time"
)

const ACCESS_TOKEN_TYPE = "ACCESS_TOKEN_TYPE"


func GenerateJWT(tokenType string , email string) (string, error) {
	var secretKey string
	var exp string
	
	if tokenType == ACCESS_TOKEN_TYPE {
		exp = time.Now().Add(12 * time.Hour).Format(time.RFC3339)
		secretKey = helpers.GetEnvValue("ACCESS_TOKEN_SECRET")
	} else {
		log.Println("Wrong token type")
		return "", errors.New("Wrong token type")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp
	claims["email"] = email

	tokenString, err := token.SignedString([]byte(secretKey));
	
	return tokenString, err
}
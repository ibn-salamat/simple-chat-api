package tools

import (
	"errors"
	"ibn-salamat/simple-chat-api/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

const ACCESS_TOKEN_TYPE = "ACCESS_TOKEN_TYPE"

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(tokenType string, email string) (string, error) {
	var secretKey string

	if tokenType == ACCESS_TOKEN_TYPE {
		secretKey = config.EnvData.ACCESS_TOKEN_SECRET
	} else {
		log.Println("Wrong token type")
		return "", errors.New("Wrong token type")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["email"] = email

	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

func CheckToken(tokenString string) error {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.EnvData.ACCESS_TOKEN_SECRET), nil
	})

	if err != nil {
		return err
	}

	return nil
}

package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"ibn-salamat/simple-chat-api/helpers"
	"log"
	"time"
)

const ACCESS_TOKEN_TYPE = "ACCESS_TOKEN_TYPE"

type Claims struct{
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(tokenType string , email string) (string, error) {
	var secretKey string
	
	if tokenType == ACCESS_TOKEN_TYPE {
		secretKey = helpers.GetEnvValue("ACCESS_TOKEN_SECRET")
	} else {
		log.Println("Wrong token type")
		return "", errors.New("Wrong token type")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(30 * time.Second).Unix()
	claims["email"] = email

	tokenString, err := token.SignedString([]byte(secretKey));
	
	return tokenString, err
}

func CheckToken (tokenString string) error {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(helpers.GetEnvValue("ACCESS_TOKEN_SECRET")), nil
	})

	v, _ := err.(*jwt.ValidationError)

	if err != nil {
		if v.Errors == jwt.ValidationErrorExpired &&  claims.ExpiresAt > time.Now().Unix(){
			return errors.New("token is expired")
		};

		return err
	}

	return nil
}
package tools

import (
	"github.com/golang-jwt/jwt"
	"log"
	"os"
)

const REFRESH_TOKEN_TYPE = "REFRESH_TOKEN_TYPE"
const ACCESS_TOKEN_TYPE = "ACCESS_TOKEN_TYPE"


func generateJWT(tokenType string , exp string, email string) (string, error) {
	var secretKey string
	
	if tokenType == ACCESS_TOKEN_TYPE {
		secretKey = os.Getenv("ACCESS_TOKEN_SECRET")
	} else if tokenType == REFRESH_TOKEN_TYPE {
		secretKey = os.Getenv("REFRESH_TOKEN_SECRET")
	} else {
		log.Fatal("Wrong token type")
	}

	token := jwt.New(jwt.SigningMethodEdDSA)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp
	claims["email"] = email
	
	tokenString, err := token.SignedString(secretKey)
	
	return tokenString, err
}
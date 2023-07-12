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

func CheckToken (tokenString string) error {
	var claims struct{
		Exp string `json:"exp"`
		jwt.StandardClaims
	}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(helpers.GetEnvValue("ACCESS_TOKEN_SECRET")), nil
	})

	v, _ := err.(*jwt.ValidationError)

	if err != nil {
		exp, timeErr := time.Parse(time.RFC3339, claims.Exp)
		if timeErr != nil {
			log.Println(err)
			return timeErr
		}

		if v.Errors == jwt.ValidationErrorExpired &&  exp.Unix() > time.Now().Unix() - (86400 * 14){
			return errors.New("token is expired")
		};

		return err
	}

	return nil
}
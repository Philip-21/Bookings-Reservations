package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Username string
	jwt.StandardClaims
}

func GenerateToken(username string) (signedToken string, signedRefreshToken string, err error) {

	//generate a token
	Payload := &SignedDetails{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	RefreshPayload := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	//call the jwt
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, Payload).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	RefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, RefreshPayload).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, RefreshToken, err
}

// confirms the token to be used in the middlewre
func ValidateToken(signedToken string) (Payload *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	Payload, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("invalid token ")
		msg = err.Error()
		return
	}

	if Payload.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return Payload, msg
}

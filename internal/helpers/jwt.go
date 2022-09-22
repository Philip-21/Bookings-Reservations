package helpers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateJWT(email string) (string, error) {

	var mySigningKey = []byte(SECRET_KEY)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(r *http.Request) (string, error) {
	if r.Header["Token"] != nil {
		var w http.ResponseWriter
		//retrieve the token from the header
		tokenString := r.Header["Token"][0]
		//parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//verify the signature
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				http.Error(w, "No Authorization Headers", http.StatusInternalServerError)
				http.Error(w, "Error with Signing Methos", http.StatusInternalServerError)
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return SECRET_KEY, nil
		})
		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			email := claims["email"].(string)
			return email, nil
		}
	}

	return "unable to extract claims", nil
}

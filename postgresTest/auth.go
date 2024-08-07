package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

func withJWTAuth(w http.ResponseWriter, r *http.Request, handlerFunc apiFunc) error {
	log.Println("running middleware")
	tokenString := r.Header.Get("x-jwt-token")

	token, err := validateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("token provided is invalid")
	}

	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(claims, "claims")
	return handlerFunc(w, r)
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo0NjQwLCJ0dGwiOjE3MjI5NDcwOTl9.Yxe-LJ48rKVmASmN5duTTx3TkCsw8fMbB7apvuA_p4A

func createJWT(account *Account) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ttl":           time.Now().Add(time.Hour * 3).Unix(),
		"accountNumber": account.Number,
	}) // Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}

func validateJWT(userToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return token, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["id"], claims["ttl"])
		fmt.Println("test claims ", claims)
	} else {
		fmt.Println(err, "errasdfasdfsaf")
	}

	return token, nil
}

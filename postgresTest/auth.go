package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("running middleware")
		tokenString := r.Header.Get("x-jwt-token")

		_, err := validateJWT(tokenString)
		if err != nil {
			// fmt.Println(err, "err")
			writeJSON(w, http.StatusUnauthorized, ApiError{Err: "Token provided is invalid"})
			return
		}

		// fmt.Println("token ", (*token).Header)

		handlerFunc(w, r)
	}
}

func validateJWT(userToken string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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

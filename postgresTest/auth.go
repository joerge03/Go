package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var secret = os.Getenv("JWT_SECRET")

func PermissionDenied(w http.ResponseWriter) error {
	return writeJSON(w, http.StatusBadRequest, ApiError{Err: "Permission denied"})
}

func withJWTAuth(w http.ResponseWriter, r *http.Request, s Storage, handlerFunc apiFunc) error {
	account := new(Account)
	log.Println("running middleware")
	tokenString := r.Header.Get("x-jwt-token")

	token, err := validateJWT(tokenString)
	if err != nil {
		return PermissionDenied(w)
	}

	claims := token.Claims.(jwt.MapClaims)
	var formattedAccountNumber int

	queries := mux.Vars(r)
	accountNumber, ok := queries["accountNumber"]

	if ok {
		formattedAccountNumber, err = strconv.Atoi(accountNumber)
		if err != nil {
			return err
		}

		account, err = s.getAccountByID(formattedAccountNumber)
		if err != nil {
			return err
		}
	}

	if ok && (claims["accountNumber"] != account.Number) {
		return PermissionDenied(w)
	}
	// fmt.Println(claims, "claims")
	return handlerFunc(w, r)
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo0NjQwLCJ0dGwiOjE3MjI5NDcwOTl9.Yxe-LJ48rKVmASmN5duTTx3TkCsw8fMbB7apvuA_p4A

// test unix
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjozMDQ3LCJ0dGwiOjE3MjMwMzEzMjZ9.mRswnhQvAwsusMTJLsqI41J09wvWSjUP_WwiwM4tLT0

// test 123
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo3NjE1LCJ0dGwiOjE3MjMxMDk4NDV9.lFX6juQs3OPp7Ibf_A21TP5iteQDD1cHk-KhrVZx3M8

func createJWT(account *Account) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ttl":           time.Now().Add(time.Hour).Unix(),
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
	// if claims, ok := token.Claims.(jwt.MapClaims); ok {
	// 	// fmt.Println(claims["id"], claims["ttl"])
	// 	// fmt.Println("test claims ", claims)
	// } else {
	// 	fmt.Println(err, "errasdfasdfsaf")
	// }
	return token, nil
}

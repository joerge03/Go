package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type ApiError struct {
	Error string `json:"error"`
}

func (c *Config) AuthWithJWT(w http.ResponseWriter, r *http.Request, handleFunc apiFunc) error {
	headerToken := r.Header.Get("x-jwt-token")

	_, err := c.ValidateJWT(headerToken)
	if err != nil {
		return err
	}

	return handleFunc(w, r)
}

func (c *Config) ValidateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_TOKEN")
	claims := new(jwt.Claims)

	tokenWithClaims, err := jwt.ParseWithClaims(token, *claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return tokenWithClaims, nil
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxSize := int64(1048576)

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})

	if err != io.EOF {
		return fmt.Errorf("enter only 1 file per upload")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data any, status int, headers ...http.Header) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for i, header := range headers[0] {
			w.Header()[i] = header
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	defaultErrorStatus := http.StatusBadRequest

	if len(status) > 0 && status[0] != 0 {
		defaultErrorStatus = status[0]
	}

	var payload jsonResponse

	payload.Message = err.Error()
	payload.Error = true

	return WriteJSON(w, payload, defaultErrorStatus)
}

func (c *Config) authenticate(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("authenticate\n")
	req := new(LoginPayload)
	res := new(TokenResponse)

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed")
	}
	err := json.NewDecoder(r.Body).Decode(req)
	fmt.Printf("recieved %v\n", req)
	if err != nil {
		return fmt.Errorf("error decoder")
	}

	user, err := c.Models.User.GetByEmail(req.Email)
	if err != nil {
		return err
	} else if user == nil {
		return fmt.Errorf("access denied")
	}

	err = user.PasswordMatches(req.Password)
	if err != nil {
		fmt.Printf("error pass %v", err)
		return err
	}

	token, err := user.CreateJWT(time.Minute * 3)
	if err != nil {
		fmt.Printf("token err %v", err)
		return err
	}

	res.Token = token
	res.ID = user.ID

	fmt.Printf(`return by auth - %v\n`, res)

	return WriteJSON(w, res, http.StatusAccepted)
}

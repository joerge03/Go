package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type LoginPayload struct {
	Email    string `json:email`
	Password string `json:password`
}
type LoginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
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

	if status[0] != 0 {
		defaultErrorStatus = status[0]
	}

	var payload jsonResponse

	payload.Message = err.Error()
	payload.Error = true

	return WriteJSON(w, payload, defaultErrorStatus)
}

func (c *Config) authenticate(w http.ResponseWriter, r *http.Request) (string, error) {
	secret := os.Getenv("JWT_TOKEN")

	req := new(LoginPayload)
	res := new(LoginResponse)

	if r.Method != "POST" {
		return "", fmt.Errorf("Method unavailable")
	}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return "", nil
	}

	user, err := c.Models.User.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}

	err = user.PasswordMatches(req.Password)
	if err != nil {
		return "", err
	}
	token, err := user.CreateJWT([]byte(secret))
	if err != nil {
		return "", err
	}
	res.Token = token
	res.ID = user.ID

	err = WriteJSON(w, res, http.StatusAccepted)
	if err != nil {
		return "", err
	}
	return token, nil
}

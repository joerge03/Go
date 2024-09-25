package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type loginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func (app *Config) JsonReader(w http.ResponseWriter, r *http.Request, data any) error {
	MaxByte := int64(1048576)

	r.Body = http.MaxBytesReader(w, r.Body, MaxByte)

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("please enter only one JSON")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	fmt.Println(headers)

	out, err := json.Marshal(data)
	if err != nil {
		log.Fatal(out)
	}

	if len(headers) > 0 {
		for i, header := range headers[0] {
			w.Header()[i] = header
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusAccepted {
		w.WriteHeader(status)
	}
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) != 0 {
		statusCode = status[0]
	}

	var payload JsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

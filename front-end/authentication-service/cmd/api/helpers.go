package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
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

func (c *Config) login(w http.ResponseWriter, r *http.Request) {
}

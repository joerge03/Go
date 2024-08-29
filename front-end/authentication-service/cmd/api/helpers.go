package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (c *Config) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
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

func (c *Config) WriteJSON(w http.ResponseWriter, data any, status int, headers ...http.Header) error {
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

func (c *Config) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	defaultErrorStatus := http.StatusBadRequest

	if status[0] != 0 {
		defaultErrorStatus = status[0]
	}

	var payload jsonResponse

	payload.Message = err.Error()
	payload.Error = true

	return c.WriteJSON(w, payload, defaultErrorStatus)
}

func (c *Config) login(w http.ResponseWriter, r *http.Request) {
}

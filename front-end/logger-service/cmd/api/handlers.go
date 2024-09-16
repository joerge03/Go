package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonPayload struct {
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
}

type JsonResponse struct{}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func WriteError(w http.ResponseWriter, err error, status ...http.ConnState) {
}

func routeHandler(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Fatal(err)
		}
	}
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) error {
	jsonPayload := new(JsonPayload)

	err := json.NewDecoder(r.Body).Decode(jsonPayload)
	if err != nil {
		return err
	}
	return nil
}

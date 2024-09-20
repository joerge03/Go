package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"logger-service/data"
)

type JsonPayload struct {
	Name string `bson:"name" json:"name"`
	Data any    `bson:"data" json:"data"`
}

type JsonResponse struct {
	Error   bool   `bson:"error" json:"error"`
	Data    any    `bson:"data" json:"data"`
	Message string `bson:"message" json:"message"`
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func (app *Config) WriteJson(w http.ResponseWriter, data any, statusCode int, headers ...http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf(`something wrong formatting json: %v\n`, err)
	}

	if len(headers) > 0 {
		for i, header := range headers[0] {
			w.Header()[i] = header
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(res)
	if err != nil {
		return fmt.Errorf(`there is something wrong writing res, %v`, err)
	}
	return nil
}

func (app *Config) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := int64(1048576)

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		return fmt.Errorf(`something wrong decoding body :%v \n`, err)
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf(`please upload only 1 file at time : %v\n`, err)
	}
	return nil
}

func (app *Config) WriteError(w http.ResponseWriter, errMessage error, status ...int) {
	errorRes := new(JsonResponse)

	defaultStatus := http.StatusAccepted
	if len(status) > 0 {
		defaultStatus = status[0]
	}

	errorRes.Error = true
	errorRes.Message = errMessage.Error()

	// TODO: CREATE WRITE JSON

	err := app.WriteJson(w, errorRes, defaultStatus)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *Config) routeHandler(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {

			routeHandlerError := fmt.Errorf(`route handler error: %s\n`, err.Error())
			app.WriteError(w, routeHandlerError)
		}
	}
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) error {
	jsonPayload := new(JsonPayload)

	err := app.ReadJson(w, r, jsonPayload)
	fmt.Println("json payload", jsonPayload)
	if err != nil {
		return err
	}

	payload := data.LogEntry{Name: jsonPayload.Name, Data: jsonPayload.Data, CreatedAt: time.Now()}

	fmt.Println("payload", payload)
	err = app.Models.LogEntry.InsertOne(payload)
	if err != nil {
		return err
	}

	jsonResponse := new(JsonResponse)

	jsonResponse.Error = false
	jsonResponse.Message = "saved"

	err = app.WriteJson(w, jsonResponse, http.StatusAccepted)
	if err != nil {
		return err
	}
	return nil
}

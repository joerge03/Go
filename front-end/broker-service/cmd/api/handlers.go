package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthPayload struct {
	Email    string `"json:"email"`
	Password string `"json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Test broker",
	}

	app.writeJSON(w, http.StatusOK, payload, r.Header)
}

func (app *Config) handleSubmit(w http.ResponseWriter, r *http.Request) {
	reqPayload := new(RequestPayload)

	err := app.JsonReader(w, r, reqPayload)
	if err != nil {
		app.ErrorJson(w, err, http.StatusInternalServerError)
	}

	switch reqPayload.Action {
	case "auth":
		app.authorize(w, reqPayload.Auth)

	default:
		app.ErrorJson(w, fmt.Errorf("action not available"))
	}
}

func (app *Config) authorize(w http.ResponseWriter, pay AuthPayload) {
	jsonData, err := json.MarshalIndent(pay, "", "\t")
	if err != nil {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", "https://authentication-service/login", bytes.NewReader(jsonData))
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	if res.StatusCode == http.StatusUnauthorized {
		app.ErrorJson(w, fmt.Errorf("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, fmt.Errorf("error calling auth service"))
		return
	}

	jsonRes := new(jsonResponse)
	err = json.NewDecoder(res.Body).Decode(jsonRes)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("error calling auth service"))
		return
	}
	if jsonRes.Error {
		app.ErrorJson(w, fmt.Errorf("error calling auth service"))
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, jsonRes)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("error calling auth service"))
		return
	}
}

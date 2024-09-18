package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
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
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.authorize(w, reqPayload.Auth)
	case "log":
		app.logIt(w, reqPayload.Log)
	default:
		app.ErrorJson(w, fmt.Errorf("action not available"))
	}
}

func (app *Config) logIt(w http.ResponseWriter, pay LogPayload) {
	logPayload, err := json.MarshalIndent(pay, "", "\t")
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	logUrl := "http://logger-service/"

	client := http.Client{}
	req, err := http.NewRequest("POST", logUrl, bytes.NewBuffer(logPayload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer res.Body.Close()

	response := new(JsonResponse)

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, response)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
}

//
//
//
//
//
//
//
//
//

func (app *Config) authorize(w http.ResponseWriter, pay AuthPayload) {
	jsonData, err := json.MarshalIndent(pay, "", "\t")
	if err != nil {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	request, err := http.NewRequest("POST", "http://authentication-service/login", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("req err %v", err)
		app.ErrorJson(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)
	fmt.Printf(`res : %v , err? %v  `, res, err)
	if err != nil {
		fmt.Printf("err res %v", err)
		app.ErrorJson(w, err)
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		app.ErrorJson(w, fmt.Errorf("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, fmt.Errorf("error calling auth service 1"))
		return
	}

	jsonRes := new(loginResponse)

	err = json.NewDecoder(res.Body).Decode(jsonRes)
	fmt.Println("ressss", jsonRes)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("error calling auth service 2"))
		return
	}
	// if jsonRes.Error {
	// 	app.ErrorJson(w, fmt.Errorf("error calling auth service 3"))
	// 	return
	// }
	err = app.writeJSON(w, http.StatusAccepted, jsonRes)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("error calling auth service 4"))
		return
	}
	defer res.Body.Close()
}

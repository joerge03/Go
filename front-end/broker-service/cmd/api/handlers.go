package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From     string `json:"from"`
	FromName string `json:"fromName"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
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
		fmt.Printf(`auth selected \n`)
		app.authorize(w, reqPayload.Auth)
	case "log":
		fmt.Printf(`log selected\n`)
		app.logIt(w, reqPayload.Log)
	case "mail":
		fmt.Printf(`mail selected\n`)
		app.mailIt(w, reqPayload.Mail)
	default:
		app.ErrorJson(w, fmt.Errorf("action not available"))
	}
}

func (app *Config) mailIt(w http.ResponseWriter, mailPayload MailPayload) {
	pay, err := json.MarshalIndent(mailPayload, "", "\t")
	if err != nil {
		log.Fatal(err)
		app.ErrorJson(w, err)
		return
	}

	fmt.Printf(`mail pay : %v \n`, mailPayload)

	client := &http.Client{}
	const MailUrl = "http://mail-service/"

	req, err := http.NewRequest("POST", MailUrl, bytes.NewBuffer(pay))
	if err != nil {
		log.Fatal(err)
		app.ErrorJson(w, err)
		return
	}
	req.Header.Add("Content-Type", "application/type")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		app.ErrorJson(w, err)
		return
	}
	defer res.Body.Close()

	resData := new(JsonResponse)

	err = json.NewDecoder(res.Body).Decode(resData)
	if err != nil {
		log.Fatalf(`There's something wrong decoding the mail: %v\n`, err)
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

func (app *Config) logIt(w http.ResponseWriter, pay LogPayload) {
	logPayload, err := json.MarshalIndent(pay, "", "\t")
	if err != nil {
		app.ErrorJson(w, fmt.Errorf(`there's something wrong using marshal with the data: %v`, err))
		return
	}
	logUrl := "http://logger-service"
	fmt.Println("payload log it", pay)

	client := http.Client{}
	req, err := http.NewRequest("POST", logUrl, bytes.NewBuffer(logPayload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		app.ErrorJson(w, fmt.Errorf(`error in making new request in logit, %v`, err))
		return
	}

	res, err := client.Do(req)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf(`error in client do %v`, err))
		return
	}
	defer res.Body.Close()

	response := new(JsonResponse)

	err = json.NewDecoder(res.Body).Decode(response)

	fmt.Println("response log it", response)
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

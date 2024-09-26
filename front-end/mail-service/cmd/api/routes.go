package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type funcHandler func(w http.ResponseWriter, r *http.Request) error

type MailPayload struct {
	From     string `json:"from"`
	FromName string `json:"fromName"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func (app *Config) routeHandler(f funcHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Fatalf(`something wrong with the routeHandler: %v\n`, err)
		}
	}
}

func (app *Config) routes() http.Handler {
	r := chi.NewRouter()
	handler := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"PUT", "GET", "DELETE", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(handler))
	r.Use(middleware.Heartbeat("/ping"))
	r.Post("/", app.routeHandler(app.SendSMTPMessage))
	return r
}

func (app *Config) SendSMTPMessage(w http.ResponseWriter, r *http.Request) error {
	payLoad := new(MailPayload)

	fmt.Printf(`payload %v \n`, r.Body)
	err := app.JsonReader(w, r, payLoad)
	if err != nil {
		return fmt.Errorf(`there's something wrong reading json: %v\n`, err)
	}

	message := Message{
		From:     payLoad.From,
		FromName: payLoad.FromName,
		To:       payLoad.To,
		Subject:  payLoad.Subject,
	}

	err = app.Mailer.SendSMTPMessage(message)
	if err != nil {
		return fmt.Errorf(`there's something wrong sending SMTP Message : %v\n`, err)
	}

	resPay := new(JsonResponse)
	resPay.Error = false
	resPay.Message = `mail has been successfully sent`

	err = app.writeJSON(w, http.StatusAccepted, resPay)
	if err != nil {
		app.ErrorJson(w, err)
	}
	return nil
}

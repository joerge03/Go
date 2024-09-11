package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (c *Config) handleHttpFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// switch r.Method {
		// case "POST":
		// 	err = ReadJson(w, r )
		// }

		if err = f(w, r); err != nil {
			ErrorJson(w, err)
		}
	}
}

func (c *Config) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"PUT", "POST", "DELETE", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Heartbeat("/ping"))
	r.Post("/login", c.handleHttpFunc(c.handleLogin))
	// r.Post("/create", c.handleHttpFunc(c.handleCreate))

	return r
}

func (c *Config) handleLogin(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("login teststssssss")
	return c.authenticate(w, r)
}

func (c *Config) handleCreate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

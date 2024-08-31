package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

	r.Post("/login", c.handleHttpFunc(c.handleLogin))
	r.Post("/get", c.handleHttpFunc(c.handleLogin))
	return r
}

func (c *Config) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return c.authenticate(w, r)
}

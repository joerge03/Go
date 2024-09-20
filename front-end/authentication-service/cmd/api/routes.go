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
		fmt.Printf("handle http func ------- \n")
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
	fmt.Println("auth enter")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"PUT", "POST", "DELETE", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Heartbeat("/ping"))
	r.Post("/login", c.handleHttpFunc(c.handleLogin))
	// r.Post("/create", c.handleHttpFunc(c.handleCreate))
	return r
}

func (c *Config) handleLogin(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("handle login!----------\n")
	return c.authenticate(w, r)
}

// func (c *Config) handleCreate(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello!")
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pew pew")
}

func middlewareHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Fprintf(w, "middleware run!")
	next(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/ws", wsHandler)

	n := negroni.Classic()

	n.Use(negroni.HandlerFunc(middlewareHandler))
	n.UseHandler(r)
	log.Fatal(http.ListenAndServe(":8080", n))
}

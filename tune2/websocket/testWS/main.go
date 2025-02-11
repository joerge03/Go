package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func serveHome(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "This is not home", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func main() {

	r := mux.NewRouter()

	hub := newHub()

	go hub.Run()

	r.HandleFunc("/", serveHome).Methods("GET")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
		serveWS(hub, w, r)
	}).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

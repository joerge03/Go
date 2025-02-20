package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ORIGIN SERVER %v\n", time.Now())
	fmt.Fprint(w, "origin server res\n")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleHome)

	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)


func handleMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		fmt.Println("HEEEHEEEE")
		handler.ServeHTTP(w,r)
	})
}

func main3(){
	r := mux.NewRouter()

	r.HandleFunc("/test/{id:[a-z]+}",func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Println(id)
	}).Methods("GET")
	
	r.Use(handleMiddleware)

	http.ListenAndServe(":8081", r)
}
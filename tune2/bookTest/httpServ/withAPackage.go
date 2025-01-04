package main

import (
	"net/http"

	"github.com/gorilla/mux"
)


func handleMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w,r)
	})
}

func main(){
	r := mux.NewRouter()

	
	
	r.Use(handleMiddleware)
}
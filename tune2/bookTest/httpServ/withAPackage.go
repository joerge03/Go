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

func main(){
	r := mux.NewRouter()


	// r.HandleFunc("")

		
	
	r.Use(handleMiddleware)
}
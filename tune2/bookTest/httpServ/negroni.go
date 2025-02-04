package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)


type contextKey string

const ContextKey contextKey = "username"

type BadAuth struct {
	username string
	password string
}


func (b *BadAuth) middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username != b.username || password != b.password {
		http.Error(w, "Invalid creds", http.StatusUnauthorized)
		return
	}
	
	c := context.WithValue(r.Context(), ContextKey, username)
	r = r.WithContext(c)
	
	next(w,r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	username, ok:= r.Context().Value(ContextKey).(string)
	if !ok {
		http.Error(w, "not string",http.StatusBadRequest)
		return
	}
	fmt.Printf("Hi %v\n", username)
}
// func myMiddleware2(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
// 	fmt.Println("test3")
// 	next(w,r)
// 	fmt.Println("test4")
// }


func main(){
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	n := negroni.Classic()
	auth := &BadAuth{
		username: "username",
		password: "password",
	}
	n.Use(negroni.HandlerFunc(auth.middleware))
	n.UseHandler(r)
	// n.Use(negroni.HandlerFunc(myMiddleware2))
	http.ListenAndServe(":8081", n)
}
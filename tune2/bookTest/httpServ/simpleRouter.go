package main

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct{}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request){
	path := req.URL.Path

	// fmt.Printf("Hi %v", name)
	switch path {
	case "/a":
		fmt.Println("you picked a")
	case "/b": 
		fmt.Println("you picked b")
	case "/c":
		fmt.Println("you picked c")
	default:
		fmt.Print("pls pick a-c")
	}
}


func main1(){
	// http.HandleFunc("/test",handleHello)	
	router := new(Router)	
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
	}
}
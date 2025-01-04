package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	Inner http.Handler
}

func (l Logger) ServeHTTP(res http.ResponseWriter, req *http.Request){
	log.Printf("start: %s\n", time.Now().String())
	l.Inner.ServeHTTP(res,req)
	log.Printf("finish %s", time.Now().String() )
}

func Serve(res http.ResponseWriter, req *http.Request){
	fmt.Println("enlosdfa")
}

func main2(){
	// handler :=  
	logger := Logger{Inner: http.HandlerFunc(Serve)}
	http.ListenAndServe(":8080",logger )
}
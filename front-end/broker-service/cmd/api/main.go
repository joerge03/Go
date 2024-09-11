package main

import (
	"fmt"
	"log"
	"net/http"
)

const route = "8085"

type Config struct{}

func main() {
	app := Config{}
	fmt.Printf("Running on localhost:%v \n", route)
	fmt.Printf("test%v \n", route)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", route),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

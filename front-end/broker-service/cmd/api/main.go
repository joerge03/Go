package main

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

const route = "80"

type Config struct {
	RabbitConn *amqp.Connection
}

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

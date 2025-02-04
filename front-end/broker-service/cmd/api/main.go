package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	route     = "80"
	rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
)

type Config struct {
	RabbitConn *amqp.Connection
}

func main() {
	// Establish connection to RABBITMQ
	conn := connect()
	defer conn.Close()

	app := Config{
		RabbitConn: conn,
	}

	consumer, err := events.NewConsumer(conn)
	if err != nil {
		FailOnError(err, "there's something wrong creating a new consumer")
	}

	fmt.Printf("Running on localhost:%v \n", route)
	fmt.Printf("test%v \n", route)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", route),
		Handler: app.routes(),
	}
	
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
	err = consumer.Listen([]string{"test", "asdfsadf", "log.INFO"})
	FailOnError(err, "there's something wrong ")
}

func connect() *amqp.Connection {
	var count float64
	var delay time.Duration
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			count++
			fmt.Println("error :", err)
		} else {
			fmt.Println("broker rabbitMQ connected!")
			connection = conn
			break
		}

		if count >= 5 {
			FailOnError(err, "can't connect to amqp server")
			return nil
		}

		delay = time.Duration(math.Pow(count, 2)) * time.Second
		fmt.Printf(`waiting for the amqp to dial in, retrying... x%v \n`, count)
		time.Sleep(delay)
	}
	return connection
}

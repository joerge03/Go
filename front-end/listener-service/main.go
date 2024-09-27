package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf(`error: %v , message: %v`, err, message)
	}
	fmt.Println("printed out but there's no error")
}

func main() {
	const rabbitmqURL = "amqp://guest:guest@localhost:5672/"

	// CONNECT
	conn, err := amqp.Dial(rabbitmqURL)
	failOnError(err, "there's something wrong dialing amqp")
	defer conn.Close()

	// OPEN A CHANNEL

	channel, err := conn.Channel()
	failOnError(err, `There's something wrong opening a channel`)

	queue, err := channel.QueueDeclare(
		"queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, `there's something wrong with declaring queue`)

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, `there's somethings wrong with declaring queue`)
}

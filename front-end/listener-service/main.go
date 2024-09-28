package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "queue"
)

func FailOnError(err error, message string) {
	if err != nil {
		log.Fatalf(`error: %v , message: %v`, err, message)
	}
}

func main() {
	// CONNECT
	conn := connect()
	defer conn.Close()

	// OPEN A CHANNEL
	channel, err := conn.Channel()
	FailOnError(err, `There's something wrong opening a channel`)

	// DECLARE

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, `there's something wrong with declaring queue`)

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Asdf"),
		},
	)

	FailOnError(err, `something wrong with the channel `)

	// Consume

	// messages, err := channel.Consume(
	// 	queue.Name,
	// 	"",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// FailOnError(err, `there's somethings wrong with declaring queue`)
}

func connect() *amqp.Connection {
	var counts int64
	var delay time.Duration
	connection := new(amqp.Connection)

	for {
		conn, err := amqp.Dial(rabbitmqURL)
		if err != nil {
			fmt.Printf(`failed to connect x%v, retrying... `, counts)
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			FailOnError(err, `There's something wrong connecting`)
			return nil
		}
		delay = time.Second * time.Duration(math.Pow(float64(counts), 2))
		time.Sleep(delay)
	}
	return connection
}

package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"log",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, `something wrong with exchange declare`)
}

func DeclareRandomQueue(ch *amqp.Channel) amqp.Queue {
	queue, err := ch.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, `There's something wrong connecting to queue`)

	return queue
}

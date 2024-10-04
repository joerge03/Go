package events

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareExchange(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
		"log",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeclareRandomQueue(channel *amqp.Channel) amqp.Queue {
	queue, err := channel.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.Fatalf(`there's something wrong with queue declare %v`, err)
	}
	return queue
}

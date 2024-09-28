package events

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func FailOnError(err error, message string) {
	if err != nil {
		log.Fatalf(`error: %v , message: %v`, err, message)
	}
}

func NewConsumer(conn *amqp.Connection) Consumer {
	Consumer := Consumer{
		conn: conn,
	}

	Consumer.setup()

	return Consumer
}

func (c *Consumer) setup() {
	channel, err := c.conn.Channel()
	FailOnError(err, `There's something wrong connecting to channel`)

	// DECLARE

	// queue, err := channel.QueueDeclare(
	// 	"queue",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// FailOnError(err, `There's something wrong connecting to queue`)

	DeclareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type authResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
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

func (c *Consumer) Listen(topics []string) {
	channel, err := c.conn.Channel()
	FailOnError(err, `there's something wrong connecting to the channel`)
	defer channel.Close()

	queue := DeclareRandomQueue(channel)

	for _, str := range topics {
		err := channel.QueueBind(
			queue.Name,
			str,
			"log",
			false,
			nil,
		)
		FailOnError(err, `failed to consume the topic`)
	}

	messages, err := channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, `Something wrong with consume`)

	// INFINITE
	infinite := make(chan bool)

	go func() {
		for d := range messages {
			data := new(Payload)

			_ = json.NewDecoder(bytes.NewReader(d.Body)).Decode(data)
			// _ = json.Unmarshal(d.Body, data)

			go handlePayload(*data)
		}
	}()
	<-infinite
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
	}
}

func logEvent(pay Payload) {
	logPayload, err := json.MarshalIndent(pay, "", "\t")

	FailOnError(err, "theres something wrong using marshal with the data")

	logUrl := "http://authentication-service"
	fmt.Println("payload log it", pay)

	client := http.Client{}
	req, err := http.NewRequest("POST", logUrl, bytes.NewBuffer(logPayload))
	req.Header.Set("Content-Type", "application/json")

	FailOnError(err, `error in making new request in logit`)

	res, err := client.Do(req)

	FailOnError(err, `error in client do %v`)

	defer res.Body.Close()

	response := new(authResponse)

	err = json.NewDecoder(res.Body).Decode(response)

	fmt.Println("response log it", response)

	FailOnError(err, `there's something wrong decoding`)
}

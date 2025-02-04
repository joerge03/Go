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
	conn *amqp.Connection
	name string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.Setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) Setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	err = DeclareExchange(channel)
	if err != nil {
		return err
	}
	return nil
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type Response struct{}

func (c *Consumer) Listen(topics []string) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue := DeclareRandomQueue(channel)
	fmt.Println("random queue name", queue.Name)
	for _, str := range topics {
		err := channel.QueueBind(
			queue.Name,
			str,
			"log",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	infinite := make(chan bool)
	go func() {
		for d := range messages {
			payload := new(Payload)

			_ = json.NewDecoder(bytes.NewReader(d.Body)).Decode(payload)

			err := c.payloadHandler(*payload)
			if err != nil {
				log.Fatalf("there's something wrong with the payload handler %v", err)
			}
		}
	}()

	fmt.Println("waiting for the message...")

	<-infinite

	fmt.Println(queue)
	return nil
}

func (c *Consumer) payloadHandler(payload Payload) error {
	switch payload.Name {
	case "log", "event":
		//
	case "auth":
		fmt.Printf("selected : auth")
	default:
		err := logEvent(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func logEvent(payload Payload) error {
	payloadByte, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	const URL = "http://logger-service/"

	req, err := http.NewRequest("POST", URL, bytes.NewReader(payloadByte))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// THIS IS NOT THE FINAL RESPONSE MIGHT CHANGE LATER
	resData := new(any)

	err = json.NewDecoder(res.Body).Decode(resData)
	if err != nil {
		return err
	}
	return nil
}

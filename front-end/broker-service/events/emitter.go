package events

import amqp "github.com/rabbitmq/amqp091-go"

type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	conn, err := e.conn.Channel()
	if err != nil {
		return err
	}

	defer conn.Close()

	err = DeclareExchange(conn)
	if err != nil {
		return err
	}
	return nil
}

func (e *Emitter) Push(event string, key string) error {
	conn, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Publish(
		"log",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{conn: conn}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emitter, nil
}

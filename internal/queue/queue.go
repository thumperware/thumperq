package queue

import (
	"github.com/streadway/amqp"
	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/publishers"
	"github.com/thumperq/thumperq/internal/reflection"
)

type IQueue interface {
	Exchange() string
	Name() string
	Bind() (*amqp.Channel, error)
	Publish(msgBytes []byte) error
}

type queue struct {
	connection connection.IConnection
	exchange   string
	name       string
}

func NewQueue(connection connection.IConnection, exchange string, name string) IQueue {
	return &queue{
		connection: connection,
		exchange:   exchange,
		name:       name,
	}
}

func (q *queue) Exchange() string {
	return q.exchange
}

func (q *queue) Name() string {
	return q.name
}

func (q *queue) Bind() (*amqp.Channel, error) {
	methodPath := reflection.MethodPath(q.Bind)
	ch, err := q.connection.RmqConnection().Channel()
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	err = ch.ExchangeDeclare(
		q.exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	_, err = ch.QueueDeclare(
		q.name, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	err = ch.QueueBind(
		q.name,     // queue name
		"",         // routing key
		q.exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return ch, nil
}

func (q *queue) Publish(msgBytes []byte) error {
	methodPath := reflection.MethodPath(q.Publish)
	rmqPublisher := publishers.NewRmqPublisher(q.connection)
	err := rmqPublisher.Publish(msgBytes, q.exchange)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	return nil
}

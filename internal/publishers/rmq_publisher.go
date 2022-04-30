package publishers

import (
	"github.com/streadway/amqp"
	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/reflection"
)

type iRmqPublisher interface {
	Publish(busMsgBytes []byte, exchange string) error
}

type rmqPublisher struct {
	connection connection.IConnection
}

func NewRmqPublisher(connection connection.IConnection) iRmqPublisher {
	return &rmqPublisher{
		connection: connection,
	}
}

func (p *rmqPublisher) Publish(busMsgBytes []byte, exchange string) error {
	methodPath := reflection.MethodPath(p.Publish)
	ch, err := p.connection.RmqConnection().Channel()
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	err = ch.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        busMsgBytes,
		})
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	return nil
}

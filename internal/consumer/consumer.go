package consumer

import (
	"github.com/streadway/amqp"
	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/queue"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/handler"
)

type consumer[T handler.IMessage] struct {
	connection      connection.IConnection
	handler         handler.IHandler[T]
	queue           queue.IQueue
	retryQueue      queue.IQueue
	errorQueue      queue.IQueue
	executeInterval int
}

func NewConsumer[T handler.IMessage](connection connection.IConnection, handler handler.IHandler[T], queue queue.IQueue, retryQueue queue.IQueue, errorQueue queue.IQueue, executeInterval int) *consumer[T] {
	consumer := consumer[T]{
		connection:      connection,
		handler:         handler,
		queue:           queue,
		retryQueue:      retryQueue,
		errorQueue:      errorQueue,
		executeInterval: executeInterval,
	}
	return &consumer
}

func (c *consumer[T]) Consume() error {
	methodPath := reflection.MethodPath(c.Consume)
	ch, err := c.queue.Bind()
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	deliveries, err := ch.Consume(
		c.queue.Name(), // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	c.executeCommands(deliveries)
	return nil
}

func (c *consumer[T]) Queue() queue.IQueue {
	return c.queue
}

func (c *consumer[T]) executeCommands(deliveries <-chan amqp.Delivery) {
	go func(deliveries <-chan amqp.Delivery) {
		for delivery := range deliveries {
			handlerCommand := NewHandlerCommand(delivery, c.retryQueue, c.errorQueue, c.handler, c.executeInterval)
			handlerCommand.Execute()
		}
	}(deliveries)
}

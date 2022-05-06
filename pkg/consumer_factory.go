package thumperq

import (
	"fmt"

	"github.com/thumperq/thumperq/internal/consumer"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/queue"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/handler"
)

func CreateConsumer[T handler.IMessage](bus IBus, handler handler.IHandler[T]) {
	methodPath := reflection.MethodPath(CreateConsumer[T])
	handlerPath := reflection.ObjectTypePath(handler)
	handlerErrorName := fmt.Sprintf("%s_error", handlerPath)
	errQueue := queue.NewQueue(bus.Connection(), handlerErrorName, handlerErrorName)
	_, err := errQueue.Bind()
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
	msgPath := reflection.TypePath[T]()
	if bus.Config().BusConfig().RetryCount <= 0 {
		firstQueue := queue.NewQueue(bus.Connection(), msgPath, handlerPath)
		firstConsumer := consumer.NewConsumer(bus.Connection(), handler, firstQueue, nil, errQueue, 0)
		err = firstConsumer.Consume()
		if err != nil {
			panic(formatter.FormatErr(methodPath, err))
		}
		return
	}
	lastRetryName := fmt.Sprintf("%s_retry%d", handlerPath, bus.Config().BusConfig().RetryCount)
	lastRetryQueue := queue.NewQueue(bus.Connection(), lastRetryName, lastRetryName)
	lastRetryConsumer := consumer.NewConsumer(bus.Connection(), handler, lastRetryQueue, nil, errQueue, bus.Config().BusConfig().RetryIntervalMs)
	err = lastRetryConsumer.Consume()
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
	nextRetryQueue := lastRetryQueue
	for i := bus.Config().BusConfig().RetryCount - 1; i >= 1; i-- {
		retryName := fmt.Sprintf("%s_retry%d", handlerPath, i)
		retryQueue := queue.NewQueue(bus.Connection(), retryName, retryName)
		retryConsumer := consumer.NewConsumer(bus.Connection(), handler, retryQueue, nextRetryQueue, errQueue, bus.Config().BusConfig().RetryIntervalMs)
		err = retryConsumer.Consume()
		if err != nil {
			panic(formatter.FormatErr(methodPath, err))
		}
		nextRetryQueue = retryQueue
		if i == 1 {
			firstQueue := queue.NewQueue(bus.Connection(), msgPath, handlerPath)
			firstConsumer := consumer.NewConsumer(bus.Connection(), handler, firstQueue, retryQueue, errQueue, 0)
			err = firstConsumer.Consume()
			if err != nil {
				panic(formatter.FormatErr(methodPath, err))
			}
		}
	}
}

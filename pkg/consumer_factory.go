package thumperq

import (
	"fmt"
	"reflect"

	"github.com/thumperq/thumperq/internal/consumer"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/queue"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/handler"
)

func CreateConsumer[T handler.IMessage](ibus IBus, handler handler.IHandler[T]) {
	if reflect.TypeOf(ibus) != reflect.TypeOf(&bus{}) {
		return
	}
	methodPath := reflection.MethodPath(CreateConsumer[T])
	handlerPath := reflection.ObjectTypePath(handler)
	handlerErrorName := fmt.Sprintf("%s_error", handlerPath)
	errQueue := queue.NewQueue(ibus.Connection(), handlerErrorName, handlerErrorName)
	_, err := errQueue.Bind()
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
	msgPath := reflection.TypePath[T]()
	if ibus.Config().BusConfig().RetryCount <= 0 {
		firstQueue := queue.NewQueue(ibus.Connection(), msgPath, handlerPath)
		firstConsumer := consumer.NewConsumer(ibus.Connection(), handler, firstQueue, nil, errQueue, 0)
		err = firstConsumer.Consume()
		if err != nil {
			panic(formatter.FormatErr(methodPath, err))
		}
		return
	}
	lastRetryName := fmt.Sprintf("%s_retry%d", handlerPath, ibus.Config().BusConfig().RetryCount)
	lastRetryQueue := queue.NewQueue(ibus.Connection(), lastRetryName, lastRetryName)
	lastRetryConsumer := consumer.NewConsumer(ibus.Connection(), handler, lastRetryQueue, nil, errQueue, ibus.Config().BusConfig().RetryIntervalMs)
	err = lastRetryConsumer.Consume()
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
	nextRetryQueue := lastRetryQueue
	for i := ibus.Config().BusConfig().RetryCount - 1; i >= 1; i-- {
		retryName := fmt.Sprintf("%s_retry%d", handlerPath, i)
		retryQueue := queue.NewQueue(ibus.Connection(), retryName, retryName)
		retryConsumer := consumer.NewConsumer(ibus.Connection(), handler, retryQueue, nextRetryQueue, errQueue, ibus.Config().BusConfig().RetryIntervalMs)
		err = retryConsumer.Consume()
		if err != nil {
			panic(formatter.FormatErr(methodPath, err))
		}
		nextRetryQueue = retryQueue
	}
	firstQueue := queue.NewQueue(ibus.Connection(), msgPath, handlerPath)
	firstConsumer := consumer.NewConsumer(ibus.Connection(), handler, firstQueue, nextRetryQueue, errQueue, 0)
	err = firstConsumer.Consume()
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
}

package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/thumperq/thumperq/internal/busmessages"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/queue"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/handler"
)

type iHandlerCommand interface {
	Execute()
}

type handlerCommand[T handler.IMessage] struct {
	delivery        amqp.Delivery
	retryQueue      queue.IQueue
	errorQeueu      queue.IQueue
	handler         handler.IHandler[T]
	executeInterval int
}

func NewHandlerCommand[T handler.IMessage](delivery amqp.Delivery, retryQueue queue.IQueue, errorQueue queue.IQueue, handler handler.IHandler[T], executeInterval int) iHandlerCommand {
	return &handlerCommand[T]{
		delivery:        delivery,
		retryQueue:      retryQueue,
		errorQeueu:      errorQueue,
		handler:         handler,
		executeInterval: executeInterval,
	}
}

func (h *handlerCommand[T]) Execute() {
	go func() {
		busMsg, err := busmessages.NewBusMessageFromBytes(h.delivery.Body)
		if err != nil {
			h.markError(h.delivery.Body, err)
			return
		}
		msg := reflection.CreateInstance[T]()
		err = json.Unmarshal([]byte(busMsg.Message), msg)
		if err != nil {
			h.markError(h.delivery.Body, err)
			return
		}
		consumerMsgStream := make(chan handler.HandlerMessage[T])
		consumerMsg := handler.NewHandlerMessage(busMsg.Headers, msg)
		go func() {
			time.AfterFunc(time.Duration(h.executeInterval*int(time.Millisecond)), func() {
				err = h.handler.Handle(consumerMsgStream)
				if err != nil {
					if h.retryQueue != nil {
						h.retry(h.delivery.Body)
					} else {
						h.markError(h.delivery.Body, err)
						go func() {
							h.handler.Compensate(consumerMsgStream)
						}()
					}
				}
			})
		}()
		consumerMsgStream <- consumerMsg
	}()
}

func (h *handlerCommand[T]) markError(msgBytes []byte, err error) {
	methodPath := reflection.MethodPath(h.markError)
	go func(msgBytes []byte, err error) {
		busErrMsg := busmessages.NewBusError(msgBytes, err)
		busErrMsgBytes, err := busErrMsg.ToJsonBytes()
		if err != nil {
			log.Default().Println(formatter.FormatErr(methodPath, err))
			return
		}
		err = h.errorQeueu.Publish(busErrMsgBytes)
		if err != nil {
			log.Default().Println(formatter.FormatErr(methodPath, err))
			return
		}
	}(msgBytes, err)
}

func (h *handlerCommand[T]) retry(msgBytes []byte) {
	methodPath := reflection.MethodPath(h.retry)
	go func(msgBytes []byte) {
		err := h.retryQueue.Publish(msgBytes)
		if err != nil {
			log.Default().Println(formatter.FormatErr(methodPath, err))
			return
		}
	}(msgBytes)
}

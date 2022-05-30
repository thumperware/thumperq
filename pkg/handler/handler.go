package handler

const HandleMethodName = "Handle"

type IMessage interface {
	Id() string
}

type IHandler[T IMessage] interface {
	Handle(msg <-chan HandlerMessage[T]) error
	Compensate(msg <-chan HandlerMessage[T])
}

type HandlerMessage[T IMessage] struct {
	Headers map[string][]string
	Message T
}

func NewHandlerMessage[T IMessage](headers map[string][]string, msg T) HandlerMessage[T] {
	return HandlerMessage[T]{
		Headers: headers,
		Message: msg,
	}
}

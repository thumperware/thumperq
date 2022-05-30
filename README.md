![alt text](https://github.com/thumperq/thumperq/blob/master/ThumperQ_Logo.png?raw=true)

Thumperq is a light weight and easy to use Go library for RabbitMQ.

# Rquirements
- Thumperq requires Go 1.8+

# How to install
```
go get github.com/thumperq/thumperq
or
go get github.com/thumperq/thumperq@v1.1.0
```

# Quick Start
1. Define your message(aka event)
```go
type MyEvent struct {
	CorrelationID                   string    `json:"correlationId,omitempty"`
	Description          string    `json:"description,omitempty"`
}

func NewMyEvent(correlationId string, description string) *MyEvent{
    return &MyEvent{
        CorrelationID: correlationId,
        Description: description,
    }
}
```
2. Create a config instance
```go
import (
    thumperqcfg "github.com/thumperq/thumperq/pkg/config"
)

busCfg := thumperqcfg.BusConfig{
	RmqConnection: "<rmq_connection_string>", // RabbitMQ connection string
	PropagateContextMetadata: true, // If set to true it'll propagate go contexts metadata in the Bus message
	RetryCount: 3, // How many retries in case of failure
	RetryIntervalMs: 2000 // 2 seconds delay in retry executing handler in case of failure
}
```
3. Create a bus instance
```go
import (
    thumperq "github.com/thumperq/thumperq/pkg"
)

bus := thumperq.NewBus(busCfg)
```
4. Publish a message(aka event)
```go
myEvent := NewMyEvent("e783a086-7d81-4e4b-bd72-325103735bfa", "Some description!")
err := bus.Publish(ctx, myEvent)
```
5. Subscribe to a message(aka event)
```go
import (
    thumperq "github.com/thumperq/thumperq/pkg"
	"github.com/thumperq/thumperq/pkg/handler"
)

type MyHandler struct {
	bus thumperq.IBus
}

func NewMyHandler(bus thumperq.IBus) *MyHandler {
	myHandler := &MyHandler{
		bus: bus,
	}
	thumperq.CreateConsumer[*MyEvent](bus, myHandler)
	return myHandler
}

func (h *MyHandler) Handle(msg <-chan handler.HandlerMessage[*MyEvent]) error {
    busMsg := <-msg
	// busMsg.Message -> is type of your event(in this example it's of type *MyEvent)
	// busMsg.Headers -> contains contexts' metadata published by the bus if the PropagateContextMetadata in config is set to true
	...
}
```
6. Compensate if all retry fails
```go
func (h *MyHandler) Compensate(msg <-chan handler.HandlerMessage[*MyEvent]) {
    busMsg := <-msg
	// busMsg.Message -> is type of your event(in this example it's of type *MyEvent)
	// busMsg.Headers -> contains contexts' metadata published by the bus if the PropagateContextMetadata in config is set to true
	// this is usefule to publish a message to notify other services to compensate the the call. the event correlation id can be used as to track messages in other services.
	...
}
```
![alt text](https://github.com/thumperq/thumperq/blob/master/ThumperQ_Logo.png?raw=true)

Thumperq is a light weight and easy to use Go library for RabbitMQ.

# Rquirements
- Thumperq requires Go 1.8+

# How to install
```
go get github.com/thumperq/thumperq
or
go get github.com/thumperq/thumperq@v1.0.0
```

# Quick Start
1. Define your message(aka event)
```go
type MyEvent struct {
	ID                   string    `json:"id,omitempty"`
	Description          string    `json:"description,omitempty"`
}

func NewMyEvent(id string, description string) *MyEvent{
    return &MyEvent{
        ID: id,
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
myEvent := NewMyEvent("1234", "Some description!")
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

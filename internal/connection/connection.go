package connection

import (
	"github.com/streadway/amqp"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/config"
)

type IConnection interface {
	RmqConnection() *amqp.Connection
}

type connection struct {
	cfg           config.IConfig
	rmqConnection *amqp.Connection
}

func NewConnection(cfg config.IConfig) (IConnection, error) {
	methodPath := reflection.MethodPath(NewConnection)
	conn, err := amqp.Dial(cfg.BusConfig().RmqConnection)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return &connection{
		cfg:           cfg,
		rmqConnection: conn,
	}, nil
}

func (c *connection) RmqConnection() *amqp.Connection {
	return c.rmqConnection
}

package thumperq

import (
	"context"

	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/publishers"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/config"
	"github.com/thumperq/thumperq/pkg/handler"
)

type IBus interface {
	Connection() connection.IConnection
	Config() config.IConfig
	Publish(ctx context.Context, message handler.IMessage) error
}

type bus struct {
	connection connection.IConnection
	cfg        config.IConfig
	publisher  publishers.IBusPublisher
}

func NewBus(cfg config.IConfig) IBus {
	methodPath := reflection.MethodPath(NewBus)
	conn, err := connection.NewConnection(cfg)
	if err != nil {
		panic(formatter.FormatErr(methodPath, err))
	}
	return &bus{
		connection: conn,
		cfg:        cfg,
		publisher:  publishers.NewBusPublisher(conn, cfg),
	}
}

func (b *bus) Publish(ctx context.Context, message handler.IMessage) error {
	methodPath := reflection.MethodPath(b.Publish)
	err := b.publisher.Publish(ctx, message)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	return nil
}

func (b *bus) Connection() connection.IConnection {
	return b.connection
}

func (b *bus) Config() config.IConfig {
	return b.cfg
}

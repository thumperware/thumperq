package publishers

import (
	"context"

	"github.com/thumperq/thumperq/internal/busmessages"
	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/config"
	"github.com/thumperq/thumperq/pkg/handler"
)

type IBusPublisher interface {
	Publish(ctx context.Context, message handler.IMessage) error
}

type busPublisher struct {
	connection connection.IConnection
	cfg        config.IConfig
}

func NewBusPublisher(connection connection.IConnection, cfg config.IConfig) IBusPublisher {
	return &busPublisher{
		connection: connection,
		cfg:        cfg,
	}
}

func (b *busPublisher) Publish(ctx context.Context, msg handler.IMessage) error {
	methodPath := reflection.MethodPath(b.Publish)
	busMsg, err := busmessages.NewBusMessage(ctx, b.cfg, msg)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	busMsgBytes, err := busMsg.ToJsonBytes()
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	rmqPublisher := NewRmqPublisher(b.connection)
	err = rmqPublisher.Publish(busMsgBytes, busMsg.MessageType)
	if err != nil {
		return formatter.FormatErr(methodPath, err)
	}
	return nil
}

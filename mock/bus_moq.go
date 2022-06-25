package mock

import (
	"context"

	"github.com/thumperq/thumperq/internal/connection"
	"github.com/thumperq/thumperq/internal/publishers"
	thumperq "github.com/thumperq/thumperq/pkg"
	"github.com/thumperq/thumperq/pkg/config"
	"github.com/thumperq/thumperq/pkg/handler"
)

type busMoq struct {
	connection connection.IConnection
	cfg        config.IConfig
	publisher  publishers.IBusPublisher
}

func NewBusMoq(cfg config.IConfig) thumperq.IBus {
	return &busMoq{
		connection: nil,
		cfg:        cfg,
		publisher:  nil,
	}
}

func (b *busMoq) Publish(ctx context.Context, message handler.IMessage) error {
	return nil
}

func (b *busMoq) Connection() connection.IConnection {
	return b.connection
}

func (b *busMoq) Config() config.IConfig {
	return b.cfg
}

package busmessages

import (
	"context"
	"encoding/json"
	"time"

	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/reflection"
	"github.com/thumperq/thumperq/pkg/config"
	"github.com/thumperq/thumperq/pkg/handler"
	"google.golang.org/grpc/metadata"
)

type BusMessage struct {
	Headers         map[string][]string
	MessageType     string
	CreatedDateTime time.Time
	Message         string
}

func NewBusMessage(ctx context.Context, cfg config.IConfig, msg handler.IMessage) (*BusMessage, error) {
	methodPath := reflection.MethodPath(NewBusMessage)
	msgJsonBytes, err := json.Marshal(msg)
	msgPath := reflection.ObjectTypePath(msg)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return &BusMessage{
		Headers:         loadHeaders(ctx, cfg),
		MessageType:     msgPath,
		CreatedDateTime: time.Now().UTC(),
		Message:         string(msgJsonBytes),
	}, nil
}

func NewBusMessageFromBytes(busMsgBytes []byte) (*BusMessage, error) {
	methodPath := reflection.MethodPath(NewBusMessageFromBytes)
	busMsg := &BusMessage{}
	err := json.Unmarshal(busMsgBytes, busMsg)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return busMsg, nil
}

func loadHeaders(ctx context.Context, cfg config.IConfig) map[string][]string {
	if cfg.BusConfig().PropagateContextMetadata {
		if md, exist := metadata.FromIncomingContext(ctx); exist {
			return md.Copy()
		}
	}
	return nil
}

func (b *BusMessage) ToJsonBytes() ([]byte, error) {
	methodPath := reflection.MethodPath(b.ToJsonBytes)
	messageJsonBytes, err := json.Marshal(b)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return messageJsonBytes, nil
}

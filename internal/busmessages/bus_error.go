package busmessages

import (
	"encoding/json"

	"github.com/thumperq/thumperq/internal/formatter"
	"github.com/thumperq/thumperq/internal/reflection"
)

type BusError struct {
	Message string
	Error   string
}

func NewBusError(msgBytes []byte, err error) *BusError {
	return &BusError{
		Message: string(msgBytes),
		Error:   err.Error(),
	}
}

func (b *BusError) ToJsonBytes() ([]byte, error) {
	methodPath := reflection.MethodPath(b.ToJsonBytes)
	messageJsonBytes, err := json.Marshal(b)
	if err != nil {
		return nil, formatter.FormatErr(methodPath, err)
	}
	return messageJsonBytes, nil
}

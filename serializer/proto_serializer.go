package serializer

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrNotImplementProtoMessage refers to param not implemented by proto.Message
var ErrNotImplementProtoMessage = errors.New("param does not implement proto.Message")

var Proto = ProtoSerializer{}

// ProtoSerializer implements the Serializer interface
type ProtoSerializer struct{}

// Marshal .
func (ProtoSerializer) Marshal(message any) ([]byte, error) {
	if message == nil {
		return []byte{}, nil
	}

	body, ok := message.(proto.Message)
	if ok {
		return proto.Marshal(body)
	}

	return nil, ErrNotImplementProtoMessage
}

// Unmarshal .
func (ProtoSerializer) Unmarshal(data []byte, message any) error {
	if message == nil {
		return nil
	}

	body, ok := message.(proto.Message)
	if ok {
		return proto.Unmarshal(data, body)
	}

	return ErrNotImplementProtoMessage
}

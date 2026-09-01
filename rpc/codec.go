package rpc

import (
	"bytes"
	"fmt"
)

const MessageTypeSize = 1
const MaxPayloadSize = 1 << 20 // 1MiB = 1024KiB = 1024 * 1024B = 2^20 B

func Encode(message Message) ([]byte, error) {
	var buf bytes.Buffer

	if err := buf.WriteByte(byte(message.MsgType())); err != nil {
		return nil, err
	}

	switch m := message.(type) {
	case *ProposeRequest:
		if _, err := buf.Write(m.Data); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported message type %T", message)
	}

	payload := buf.Bytes()
	if len(payload) > MessageTypeSize+MaxPayloadSize {
		return nil, fmt.Errorf("payload too large")
	}
	return payload, nil
}

func Decode(data []byte) (Message, error) {
	if len(data) < MessageTypeSize {
		return nil, fmt.Errorf("payload too small")
	}
	if len(data) > MessageTypeSize+MaxPayloadSize {
		return nil, fmt.Errorf("payload too large")
	}

	msgType := MessageType(data[0])
	payload := data[MessageTypeSize:]

	switch msgType {
	case TypeProposeRequest:
		return &ProposeRequest{
			Data: bytes.Clone(payload), // It's safer to clone bytes here, because a caller may change the original buffer. The slice may be changed after this instance created.
		}, nil
	default:
		return nil, fmt.Errorf("unknown message type %T", msgType)
	}
}

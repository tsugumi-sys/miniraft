package rpc

import (
	"bytes"
	"fmt"
)

const MaxMessageSize = 1 << 20 // 1MiB = 1024KiB = 1024 * 1024B = 2^20 B

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
	if len(payload) > MaxMessageSize {
		return nil, fmt.Errorf("payload too large")
	}
	return payload, nil
}

// func Decode(payload []byte) (Message error) {
// }

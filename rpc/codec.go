package rpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Encode(message Message) (Frame, error) {
	if !validateMessageType(message) {
		return Frame{}, fmt.Errorf("invalid message type: %d", message.MsgType())
	}
	payload, err := message.EncodePayload()
	if err != nil {
		return Frame{}, fmt.Errorf("failed to encode payload: %w", err)
	}
	payloadSize := len(payload)
	if payloadSize > MaxPayloadSize {
		return Frame{}, fmt.Errorf("payload too large, got %d bytes, maximum is %d", payloadSize, MaxPayloadSize)
	}

	frame := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	binary.BigEndian.PutUint32(frame[:FrameHeaderSize], uint32(payloadSize))
	frame[FrameHeaderSize] = byte(message.MsgType())
	copy(frame[FrameHeaderSize+MessageTypeSize:], payload)
	return Frame{data: frame}, nil
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

package rpc

import (
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

func Decode(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}

	msgType := MessageType(frame.data[FrameHeaderSize])
	decoder, ok := decoders[msgType]
	if !ok {
		return nil, fmt.Errorf("decoder not supported for the message type: %d", msgType)
	}
	return decoder(frame)
}

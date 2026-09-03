package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Frame struct {
	// [4-byte payload length][1-byte message type][N-byte payload]
	data []byte
}

const FrameHeaderSize = 4
const MessageTypeSize = 1
const MaxPayloadSize = 1 << 20 // 1MiB

func ReadFrame(r io.Reader) (Frame, error) {
	var header [FrameHeaderSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return Frame{}, fmt.Errorf("failed to read the frame header: %w", err)
	}

	payloadSize := binary.BigEndian.Uint32(header[:])
	if payloadSize > MaxPayloadSize {
		return Frame{}, fmt.Errorf("the payload too large: %d, maximum: %d", payloadSize, MaxPayloadSize)
	}

	frame := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	copy(frame[:FrameHeaderSize], header[:])
	if _, err := io.ReadFull(r, frame[FrameHeaderSize:]); err != nil {
		return Frame{}, fmt.Errorf("failed to read the frame payload: %w", err)
	}
	return Frame{data: frame}, nil

}

func WriteFrame(w io.Writer, frame Frame) error {
	if len(frame.data) < FrameHeaderSize {
		return errors.New("the frame too short")
	}
	payloadSize := binary.BigEndian.Uint32(frame.data[:FrameHeaderSize])
	if payloadSize > MaxPayloadSize {
		return fmt.Errorf("the payload too large: %d, maximum: %d", payloadSize, MaxPayloadSize)
	}

	data := frame.data
	for len(frame.data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write the frame to writer: %w", err)
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}

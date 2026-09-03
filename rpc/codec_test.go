package rpc

import (
	"encoding/binary"
	"testing"
)

func TestEncodeAcceptsMaxPayload(t *testing.T) {
	message := &ProposeRequest{
		Data: make([]byte, 1<<20),
	}
	if _, err := Encode(message); err != nil {
		t.Fatalf("Encode() unexpected error %v", err)
	}
}

func TestEncodeRejectPayloadOverLimit(t *testing.T) {
	message := &ProposeRequest{
		Data: make([]byte, 1<<20+1),
	}
	if _, err := Encode(message); err == nil {
		t.Fatalf("Encode accepts too-large payload.")
	}
}

func TestDecodeRejectTooSmallPayload(t *testing.T) {
	data := make([]byte, 0)
	if _, err := Decode(Frame{data: data}); err == nil {
		t.Fatalf("Decode accepts an empty data")
	}
}

func TestDecodeRejectPayloadOverLimit(t *testing.T) {
	payloadSize := 1<<20 + 1
	data := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	if _, err := Decode(Frame{data: data}); err == nil {
		t.Fatalf("Decode accepts an too large data")
	}
}

func TestDecodeAcceptsMaxPayload(t *testing.T) {
	payloadSize := 1 << 20
	data := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	data[FrameHeaderSize] = byte(TypeProposeRequest)
	binary.BigEndian.PutUint32(data[:FrameHeaderSize], uint32(payloadSize))
	if _, err := Decode(Frame{data: data}); err != nil {
		t.Fatalf("Decode unexpected error %v", err)
	}
}

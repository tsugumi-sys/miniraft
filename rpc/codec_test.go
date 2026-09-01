package rpc

import "testing"

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
	if _, err := Decode(data); err == nil {
		t.Fatalf("Decode accepts an empty data")
	}
}

func TestDecodeRejectPayloadOverLimit(t *testing.T) {
	data := make([]byte, 1<<20+1)
	if _, err := Decode(data); err == nil {
		t.Fatalf("Decode accepts an too large data")
	}
}

func TestDecodeAcceptsMaxPayload(t *testing.T) {
	data := make([]byte, 1<<20)
	data[0] = byte(TypeProposeRequest)
	if _, err := Decode(data); err != nil {
		t.Fatalf("Decode unexpected error %v", err)
	}
}

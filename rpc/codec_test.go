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

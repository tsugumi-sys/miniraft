package rpc

import "testing"

func TestEncodeMaxPayloadError(t *testing.T) {
	message := ProposeRequest{
		Data: make([]byte, 1<<21),
	}
}

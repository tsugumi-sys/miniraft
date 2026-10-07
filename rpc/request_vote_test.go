package rpc

import (
	"bytes"
	"testing"
)

func TestEncodeRequestVoteRequest(t *testing.T) {
	message := &VoteRequest{
		Term:           1,
		CandidateID:    2,
		LatestLogIndex: 3,
		LatestLogTerm:  4,
	}
	payload, err := message.EncodePayload()
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}
	want := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Term
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // CandidateID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // LatestLogIndex
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // LatestLogTerm
	}
	if !bytes.Equal(payload, want) {
		t.Fatalf("EncodePayload() got: %x, want: %x", payload, want)
	}
}

func TestVoteRequestMsgType(t *testing.T) {
	message := &VoteRequest{
		Term:           1,
		CandidateID:    2,
		LatestLogIndex: 3,
		LatestLogTerm:  4,
	}
	want := TypeVoteRequest
	if message.MsgType() != want {
		t.Fatalf("unexpected message type, got: %v, want: %v", message.MsgType(), want)
	}
}

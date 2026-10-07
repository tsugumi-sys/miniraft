package rpc

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncodeVoteRequest(t *testing.T) {
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

func TestDecodeVoteRequest_ShortFrame(t *testing.T) {
	frameData := make([]byte, 1) // too short frame.
	if _, err := decodeVoteRequest(Frame{data: frameData}); err == nil {
		t.Fatal("frame is not validated.")
	}
}
func TestDecodeVoteRequest_InvalidPayloadSize(t *testing.T) {
	payloadSize := 31
	frameData := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	binary.BigEndian.PutUint32(frameData[:FrameHeaderSize], uint32(payloadSize))
	frameData[FrameHeaderSize] = byte(TypeVoteRequest)
	if _, err := decodeVoteRequest(Frame{data: frameData}); err == nil {
		t.Fatal("payload size is not validated")
	}
}
func TestDecodeVoteRequest_InvalidMessageType(t *testing.T) {
	payloadSize := VoteRequestPayloadSize
	frameData := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	binary.BigEndian.PutUint32(frameData[:FrameHeaderSize], uint32(payloadSize))
	frameData[FrameHeaderSize] = byte(TypeVoteResponse)
	if _, err := decodeVoteRequest(Frame{data: frameData}); err == nil {
		t.Fatal("message type is not validated")
	}
}
func TestDecodeVoteRequest_ValidFrame(t *testing.T) {
	payloadSize := VoteRequestPayloadSize
	frameData := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	binary.BigEndian.PutUint32(frameData[:FrameHeaderSize], uint32(payloadSize))
	frameData[FrameHeaderSize] = byte(TypeVoteRequest)
	offset := FrameHeaderSize + MessageTypeSize
	for i := 1; i <= 4; i++ {
		binary.BigEndian.PutUint64(frameData[offset:offset+8], uint64(i))
		offset += 8
	}

	msg, err := decodeVoteRequest(Frame{data: frameData})
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	got, ok := msg.(*VoteRequest)
	if !ok {
		t.Fatalf("unexpected message type: %T", msg)
	}
	want := VoteRequest{
		Term:           1,
		CandidateID:    2,
		LatestLogIndex: 3,
		LatestLogTerm:  4,
	}
	if *got != want {
		t.Fatalf("decodeVoteRequest() got: %v, want: %v", got, want)
	}
}

func TestEncodeVoteResponse(t *testing.T) {
	msg := &VoteResponse{
		VoterTerm:   1,
		VoterID:     2,
		VoteGranted: true,
	}
	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("EncodePayload() got: %v", err)
	}
	want := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // VoterTerm
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // VoterID
		0x01, // VoteGranted
	}
	if !bytes.Equal(encoded, want) {
		t.Fatalf("EncodePayload() got: %x, want: %x", encoded, want)
	}
}

func TestVoteResponseMsgType(t *testing.T) {
	msg := &VoteResponse{
		VoterTerm:   1,
		VoterID:     2,
		VoteGranted: true,
	}

	want := TypeVoteResponse
	if msg.MsgType() != want {
		t.Fatalf("unexpected message type, got: %v, want: %v", msg.MsgType(), want)
	}
}

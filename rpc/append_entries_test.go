package rpc

import (
	"encoding/binary"
	"testing"
)

func TestEncodeAppendEntriesRequest(t *testing.T) {
	entries := make([]LogEntry, 1)
	entries = append(entries, LogEntry{
		Index:   1,
		Term:    1,
		Payload: []byte("hello"),
	})
	message := &AppendEntriesRequest{
		Term:         1,
		LeaderID:     1,
		PrevLogIndex: 1,
		PrevLogTerm:  1,
		LeaderCommit: 1,
		Entries:      make([]LogEntry, 0),
	}

	if _, err := message.EncodePayload(); err != nil {
		t.Fatalf("Encode failed unexpectedly: %v", err)
	}
}

func TestAppendEntriesRequestMsgType(t *testing.T) {
	message := &AppendEntriesRequest{
		Term:         1,
		LeaderID:     1,
		PrevLogIndex: 1,
		PrevLogTerm:  1,
		LeaderCommit: 1,
		Entries:      make([]LogEntry, 0),
	}
	if message.MsgType() != TypeAppendEntriesRequest {
		t.Fatalf("unexpected message type. expected: %v, got: %v", TypeAppendEntriesRequest, message.MsgType())
	}
}

func TestDecodeAppendEntriesRequest(t *testing.T) {
	payloadSize := AppendEntriesRequestMetaSize + 20
	frameData := make([]byte, FrameHeaderSize+MessageTypeSize+payloadSize)
	// Frame metadata
	binary.BigEndian.PutUint32(frameData[:FrameHeaderSize], uint32(payloadSize))
	frameData[FrameHeaderSize] = byte(TypeAppendEntriesRequest)

	// payload
	offset := FrameHeaderSize + MessageTypeSize
	for i := 1; i < 6; i++ {
		binary.BigEndian.PutUint64(frameData[offset:offset+8], uint64(i))
		offset += 8
	}
	// payload: LogEntry
	binary.BigEndian.PutUint64(frameData[offset:offset+8], 100) // Index
	offset += 8
	binary.BigEndian.PutUint64(frameData[offset:offset+8], 101) // Term
	offset += 8
	binary.BigEndian.PutUint32(frameData[offset:offset+4], 0) // payload size is zero.

	msg, err := decodeAppendEntriesRequest(Frame{data: frameData})
	if err != nil {
		t.Fatalf("failed to decode message %v", err)
	}
	req, ok := msg.(*AppendEntriesRequest)
	if !ok {
		t.Fatalf("unexpected message type: %T", msg)
	}
	if req.Term != 1 {
		t.Errorf("unexpected Term. expected: %v, but got %v", 1, req.Term)
	}
	if req.LeaderID != 2 {
		t.Errorf("unexpected LeaderID. expected: %v, but got %v", 2, req.LeaderID)
	}
	if req.PrevLogIndex != 3 {
		t.Errorf("unexpected PrevLogIndex. expected: %v, but got %v", 3, req.PrevLogIndex)
	}
	if req.PrevLogTerm != 4 {
		t.Errorf("unexpected PrevLogTerm. expected: %v, but got %v", 3, req.PrevLogTerm)
	}

	if req.LeaderCommit != 5 {
		t.Errorf("unexpected LeaderCommit. expected: %v, but got %v", 5, req.LeaderCommit)
	}
	if len(req.Entries) != 1 {
		t.Errorf("too many entries")
	}
	logEntry := req.Entries[0]
	if logEntry.Index != 100 {
		t.Errorf("unexpected logEntry.Index. expected: %v, but got %v", 100, logEntry.Index)
	}
	if logEntry.Term != 101 {
		t.Errorf("unexpected logEntry.Term. expected: %v, but got %v", 101, logEntry.Term)
	}
	if len(logEntry.Payload) != 0 {
		t.Errorf("unexpected payload size. expected: %v, but got %v", 0, len(logEntry.Payload))
	}
}

package rpc

import (
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

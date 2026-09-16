package rpc

import "encoding/binary"

type LogEntry struct {
	Index   uint64
	Term    uint64
	Payload []byte // Raft just passes the payload bytes without recognizing its payload.
}

type AppendEntriesRequest struct {
	Term         uint64
	LeaderID     uint64
	PrevLogIndex uint64
	PrevLogTerm  uint64
	LeaderCommit uint64
	Entries      []LogEntry
}

type AppendEntriesResponse struct {
	Term       uint64
	FollowerID uint64
	Success    bool
	MatchIndex uint64
}

func (m *AppendEntriesRequest) MsgType() MessageType {
	return TypeAppendEntriesRequest
}

const LogEntryMetaSize = 8*2 + 8 // Index|Term|PayloadSize
const AppendEntriesRequestMetaSize = 8 * 5

func (m *AppendEntriesRequest) EncodePayload() ([]byte, error) {
	entriesSize := 0
	for i := 0; i < len(m.Entries); i++ {
		entriesSize += LogEntryMetaSize + len(m.Entries[i].Payload)
	}
	payload := make([]byte, AppendEntriesRequestMetaSize+entriesSize)
	binary.BigEndian.PutUint64(payload[:8], m.Term)
	binary.BigEndian.PutUint64(payload[8:16], m.LeaderID)
	binary.BigEndian.PutUint64(payload[16:24], m.PrevLogIndex)
	binary.BigEndian.PutUint64(payload[24:32], m.PrevLogTerm)
	binary.BigEndian.PutUint64(payload[32:40], m.LeaderCommit)
	currentOffset := 40
	for i := 0; i < len(m.Entries); i++ {
		entry := m.Entries[i]
		binary.BigEndian.PutUint64(payload[currentOffset:currentOffset+8], entry.Index)
		currentOffset += 8
		binary.BigEndian.PutUint64(payload[currentOffset:currentOffset+8], entry.Term)
		currentOffset += 8
		entryPayloadSize := len(entry.Payload)
		binary.BigEndian.PutUint64(payload[currentOffset:currentOffset+8], uint64(entryPayloadSize))
		currentOffset += 8
		copy(payload[currentOffset:currentOffset+entryPayloadSize], entry.Payload)
	}
	return payload, nil
}

// func decodeAppendEntriesRequest(frame Frame) (Message, error) {
// }

// func (m *AppendEntriesResponse) MsgType() MessageType {
// 	return TypeAppendEntriesResponse
// }

// func decodeAppendEntriesResponse(frame Frame) (Message, error) {

// }

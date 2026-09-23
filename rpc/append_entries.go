package rpc

import (
	"encoding/binary"
	"fmt"
)

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

const AppendEntriesResponseSize = 8 * 4

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

func decodeAppendEntriesRequest(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	data := frame.data
	if len(data) < FrameHeaderSize+MessageTypeSize+AppendEntriesRequestMetaSize {
		return nil, fmt.Errorf("frame too short for AppendEntriesRequest, it should be greater than or equal to %v, but got %v", AppendEntriesRequestMetaSize, len(data))
	}
	if err := validateFrameMessageType(frame, TypeAppendEntriesRequest); err != nil {
		return nil, err
	}
	payload := data[FrameHeaderSize+MessageTypeSize:]
	var entries []LogEntry
	offset := AppendEntriesRequestMetaSize
	for offset < len(payload) {
		if len(payload[offset:]) < 20 {
			return nil, fmt.Errorf("LogEntry too short. expected > 20, but got %v", len(payload[offset:]))
		}
		index := binary.BigEndian.Uint64(payload[offset : offset+8])
		offset += 8
		term := binary.BigEndian.Uint64(payload[offset : offset+8])
		offset += 8
		entryPayloadSize := binary.BigEndian.Uint32(payload[offset : offset+4])
		offset += 4
		entries = append(entries, LogEntry{
			Index:   index,
			Term:    term,
			Payload: payload[offset : offset+int(entryPayloadSize)],
		})
		offset += int(entryPayloadSize)
	}

	return &AppendEntriesRequest{
		Term:         binary.BigEndian.Uint64(payload[:8]),
		LeaderID:     binary.BigEndian.Uint64(payload[8:16]),
		PrevLogIndex: binary.BigEndian.Uint64(payload[16:24]),
		PrevLogTerm:  binary.BigEndian.Uint64(payload[24:32]),
		LeaderCommit: binary.BigEndian.Uint64(payload[32:40]),
		Entries:      entries,
	}, nil
}

func (m *AppendEntriesResponse) MsgType() MessageType {
	return TypeAppendEntriesResponse
}

func (m *AppendEntriesResponse) EncodePayload() ([]byte, error) {
	payload := make([]byte, AppendEntriesResponseSize)
	binary.BigEndian.PutUint64(payload[:8], m.Term)
	binary.BigEndian.PutUint64(payload[8:16], m.FollowerID)
	if m.Success {
		payload[16] = 1
	} else {
		payload[16] = 0
	}
	binary.BigEndian.PutUint64(payload[:8], m.Term)
	return payload, nil
}

func decodeAppendEntriesResponse(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, AppendEntriesResponseSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeAppendEntriesResponse); err != nil {
		return nil, err
	}
	return &AppendEntriesResponse{
		Term:       binary.BigEndian.Uint64(frame.data[:8]),
		FollowerID: binary.BigEndian.Uint64(frame.data[8:16]),
		Success:    binary.BigEndian.Uint64(frame.data[16:24]) != 0,
		MatchIndex: binary.BigEndian.Uint64(frame.data[24:]),
	}, nil
}

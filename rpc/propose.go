package rpc

import (
	"encoding/binary"
	"fmt"
)

type ProposeRequest struct {
	Data []byte
}

type ProposeResponseStatus uint8

const (
	StatusUnknown ProposeResponseStatus = iota
	StatusCommitted
	StatusFailed
	NoitLeader
)

type ProposeResponse struct {
	Status   ProposeResponseStatus
	Term     uint64
	LogIndex uint64
	LeaderID uint64
}

func (s ProposeResponseStatus) Valid() bool {
	switch s {
	case StatusCommitted, StatusFailed, NoitLeader:
		return true
	default:
		return false
	}
}

func (*ProposeRequest) MsgType() MessageType {
	return TypeProposeRequest
}

func (p *ProposeRequest) EncodePayload() ([]byte, error) {
	return p.Data, nil
}

func (*ProposeResponse) MsgType() MessageType {
	return TypeProposeResponse
}

const ProposeResponsePayloadSize = 1 + 8 + 8 + 8

func (p *ProposeResponse) EncodePayload() ([]byte, error) {
	if !p.Status.Valid() {
		return nil, fmt.Errorf("invalid propose response status: %d", p.Status)
	}

	payload := make([]byte, ProposeResponsePayloadSize)
	payload[0] = byte(p.Status)
	binary.BigEndian.PutUint64(payload[1:9], p.Term)
	binary.BigEndian.PutUint64(payload[9:17], p.LogIndex)
	binary.BigEndian.PutUint64(payload[17:], p.LeaderID)
	return payload[:], nil
}

func decodeProposeRequest(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}

	f := frame.data
	messageType := MessageType(f[FrameHeaderSize])
	if messageType != TypeProposeRequest {
		return nil, fmt.Errorf("unexpected message type: got %d, want %d", messageType, TypeProposeRequest)
	}
	return &ProposeRequest{
		Data: f[FrameHeaderSize+MessageTypeSize:],
	}, nil
}

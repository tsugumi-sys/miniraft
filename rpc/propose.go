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

const ProposeResponsePayloadSize = 25

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
	if err := validateFrameMessageType(frame, TypeProposeRequest); err != nil {
		return nil, err
	}
	return &ProposeRequest{
		Data: frame.data[FrameHeaderSize+MessageTypeSize:],
	}, nil
}

func decodeProposeResponse(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, ProposeResponsePayloadSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeProposeResponse); err != nil {
		return nil, err
	}

	return &ProposeResponse{
		Status:   ProposeResponseStatus(frame.data[0]),
		Term:     binary.BigEndian.Uint64(frame.data[1:9]),
		LogIndex: binary.BigEndian.Uint64(frame.data[9:17]),
		LeaderID: binary.BigEndian.Uint64(frame.data[17:]),
	}, nil
}

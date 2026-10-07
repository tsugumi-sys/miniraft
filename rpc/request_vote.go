package rpc

import (
	"encoding/binary"
)

type VoteRequest struct {
	Term           uint64
	CandidateID    uint64
	LatestLogIndex uint64
	LatestLogTerm  uint64
}

const VoteRequestPayloadSize = 32

type VoteResponse struct {
	VoterTerm   uint64
	VoterID     uint64
	VoteGranted bool
}

const VoteResponsePayloadSize = 17

func (m *VoteRequest) MsgType() MessageType {
	return TypeVoteRequest
}

func (m *VoteRequest) EncodePayload() ([]byte, error) {
	payload := make([]byte, VoteRequestPayloadSize)
	binary.BigEndian.PutUint64(payload[:8], m.Term)
	binary.BigEndian.PutUint64(payload[8:16], m.CandidateID)
	binary.BigEndian.PutUint64(payload[16:24], m.LatestLogIndex)
	binary.BigEndian.PutUint64(payload[24:], m.LatestLogTerm)
	return payload, nil
}

func decodeVoteRequest(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, VoteRequestPayloadSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeVoteRequest); err != nil {
		return nil, err
	}
	payload, err := frame.Payload()
	if err != nil {
		return nil, err
	}
	return &VoteRequest{
		Term:           binary.BigEndian.Uint64(payload[:8]),
		CandidateID:    binary.BigEndian.Uint64(payload[8:16]),
		LatestLogIndex: binary.BigEndian.Uint64(payload[16:24]),
		LatestLogTerm:  binary.BigEndian.Uint64(payload[24:]),
	}, nil
}

func (m *VoteResponse) EncodePayload() ([]byte, error) {
	payload := make([]byte, VoteResponsePayloadSize)
	binary.BigEndian.PutUint64(payload[:8], m.VoterTerm)
	binary.BigEndian.PutUint64(payload[8:16], m.VoterID)
	if m.VoteGranted {
		payload[16] = 1
	} else {
		payload[16] = 0
	}
	return payload, nil
}

func (m *VoteResponse) MsgType() MessageType {
	return TypeVoteResponse
}

func decodeVoteResponse(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, VoteResponsePayloadSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeVoteResponse); err != nil {
		return nil, err
	}
	var voteGranted bool
	if frame.data[16] == 1 {
		voteGranted = true
	} else {
		voteGranted = false
	}
	payload, err := frame.Payload()
	if err != nil {
		return nil, err
	}
	return &VoteResponse{
		VoterTerm:   binary.BigEndian.Uint64(payload[:8]),
		VoterID:     binary.BigEndian.Uint64(payload[8:16]),
		VoteGranted: voteGranted,
	}, nil
}

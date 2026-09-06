package rpc

import (
	"encoding/binary"
)

type RequestVoteRequest struct {
	Term           uint64
	CandidateID    uint64
	LatestLogIndex uint64
	LatestLogTerm  uint64
}

const RequestVoteRequestPayloadSize = 32

type RequestVoteResponse struct {
	VoterTerm   uint64
	VoterID     uint64
	VoteGranted bool
}

const RequestVoteResponsePayloadSize = 16 + 1

func (m *RequestVoteRequest) MsgType() MessageType {
	return TypeRequestVoteRequest
}

func (m *RequestVoteRequest) EncodePayload() ([]byte, error) {
	payload := make([]byte, RequestVoteRequestPayloadSize)
	binary.BigEndian.PutUint64(payload[:8], m.Term)
	binary.BigEndian.PutUint64(payload[8:16], m.CandidateID)
	binary.BigEndian.PutUint64(payload[16:24], m.LatestLogIndex)
	binary.BigEndian.PutUint64(payload[24:], m.LatestLogTerm)
	return payload, nil
}

func decodeRequestVoteRequest(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, RequestVoteRequestPayloadSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeRequestVoteRequest); err != nil {
		return nil, err
	}
	return &RequestVoteRequest{
		Term:           binary.BigEndian.Uint64(frame.data[:8]),
		CandidateID:    binary.BigEndian.Uint64(frame.data[8:16]),
		LatestLogIndex: binary.BigEndian.Uint64(frame.data[16:24]),
		LatestLogTerm:  binary.BigEndian.Uint64(frame.data[24:]),
	}, nil
}

func (m *RequestVoteResponse) EncodePayload() ([]byte, error) {
	payload := make([]byte, RequestVoteResponsePayloadSize)
	binary.BigEndian.PutUint64(payload[8:16], m.VoterTerm)
	binary.BigEndian.PutUint64(payload[16:24], m.VoterID)
	if m.VoteGranted {
		payload[24] = 1
	} else {
		payload[24] = 0
	}
	return payload, nil
}

func (m *RequestVoteResponse) MsgType() MessageType {
	return TypeRequestVoteResponse
}

func decodeRequestVoteResponse(frame Frame) (Message, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if err := validateFramePayloadSize(frame, RequestVoteRequestPayloadSize); err != nil {
		return nil, err
	}
	if err := validateFrameMessageType(frame, TypeRequestVoteResponse); err != nil {
		return nil, err
	}
	var voteGranted bool
	if frame.data[16] == 1 {
		voteGranted = true
	} else {
		voteGranted = false
	}
	return &RequestVoteResponse{
		VoterTerm:   binary.BigEndian.Uint64(frame.data[:8]),
		VoterID:     binary.BigEndian.Uint64(frame.data[8:16]),
		VoteGranted: voteGranted,
	}, nil
}

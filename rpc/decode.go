package rpc

type MessageDecoder func(frame Frame) (Message, error)

var decoders = map[MessageType]MessageDecoder{
	TypeProposeRequest:        decodeProposeRequest,
	TypeProposeResponse:       decodeProposeResponse,
	TypeVoteRequest:           decodeVoteRequest,
	TypeVoteResponse:          decodeVoteResponse,
	TypeAppendEntriesRequest:  decodeAppendEntriesRequest,
	TypeAppendEntriesResponse: decodeAppendEntriesResponse,
}

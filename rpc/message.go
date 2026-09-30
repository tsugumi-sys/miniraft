package rpc

type MessageType uint8

const (
	TypeUnknown MessageType = iota
	TypeProposeRequest
	TypeProposeResponse
	TypeVoteRequest
	TypeVoteResponse
	TypeAppendEntriesRequest
	TypeAppendEntriesResponse
)

type Message interface {
	MsgType() MessageType
	EncodePayload() ([]byte, error)
}

func validateMessageType(message Message) bool {
	switch message.MsgType() {
	case TypeProposeRequest, TypeProposeResponse, TypeVoteRequest, TypeVoteResponse:
		return true
	default:
		return false
	}
}

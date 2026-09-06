package rpc

type MessageType uint8

const (
	TypeUnknwon MessageType = iota
	TypeProposeRequest
	TypeProposeResponse
	TypeRequestVoteRequest
	TypeRequestVoteResponse
)

type Message interface {
	MsgType() MessageType
	EncodePayload() ([]byte, error)
}

func validateMessageType(message Message) bool {
	switch message.MsgType() {
	case TypeProposeRequest, TypeProposeResponse, TypeRequestVoteRequest, TypeRequestVoteResponse:
		return true
	default:
		return false
	}
}

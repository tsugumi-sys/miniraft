package rpc

type MessageType uint8

const (
	TypeUnknwon MessageType = iota
	TypeProposeRequest
	TypeProposeResponse
)

type Message interface {
	MsgType() MessageType
	EncodePayload() ([]byte, error)
}

func validateMessageType(message Message) bool {
	switch message.MsgType() {
	case TypeProposeRequest, TypeProposeResponse:
		return true
	default:
		return false
	}
}

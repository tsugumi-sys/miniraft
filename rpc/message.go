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

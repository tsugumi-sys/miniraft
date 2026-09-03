package rpc

type MessageDecoder func(frame Frame) (Message, error)

var decoders = map[MessageType]MessageDecoder{
	TypeProposeRequest: decodeProposeRequest,
}

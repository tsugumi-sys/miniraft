package rpc

type ProposeRequest struct {
	Data []byte
}

type ProposeResponseStatus uint8

const (
	StatusUnknown ProposeResponseStatus = iota + 1
	StatusCommitted
	StatusFailed
	NoitLeader
)

type ProposeResponse struct {
	Status     ProposeResponseStatus
	Term       uint64
	LoginIndex uint64
	LeaderID   uint64
}

func (*ProposeRequest) MsgType() MessageType {
	return TypeProposeRequest
}

func (*ProposeResponse) MsgType() MessageType {
	return TypeProposeResponse
}

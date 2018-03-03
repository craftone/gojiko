package gtpv1u

import "github.com/craftone/gojiko/domain/gtp"

type EchoRequest struct {
	header
}

func NewEchoRequest(seqNum uint16) *EchoRequest {
	return &EchoRequest{
		newHeader(EchoRequestNum, gtp.Teid(0), seqNum),
	}
}

func (e *EchoRequest) Marshal() []byte {
	return e.header.marshal([]byte{})
}

type EchoResponse struct {
	header
	recovery byte
}

func NewEchoResponse(seqNum uint16, recovery byte) *EchoResponse {
	return &EchoResponse{
		newHeader(EchoRequestNum, gtp.Teid(0), seqNum),
		recovery,
	}
}

func (e *EchoResponse) Marshal() []byte {
	return e.header.marshal([]byte{14, e.recovery})
}

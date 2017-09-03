package gtpv2c

import (
	"log"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type echo struct {
	header   *header
	recovery *ie.Recovery
}

type EchoRequest struct {
	echo
}

type EchoResponse struct {
	echo
}

func newEcho(messageType messageType, seqNum uint32, recovery byte) echo {
	return echo{
		newHeader(messageType, false, false, 0, seqNum),
		ie.NewRecovery(recovery, 0),
	}
}

func NewEchoRequest(seqNum uint32, recovery byte) *EchoRequest {
	return &EchoRequest{
		newEcho(echoRequest, seqNum, recovery),
	}
}

func NewEchoResponse(seqNum uint32, recovery byte) *EchoRequest {
	return &EchoRequest{
		newEcho(echoResponse, seqNum, recovery),
	}
}

func (e *echo) Marshal() []byte {
	body := e.recovery.Marshal()
	return e.header.marshal(body)
}

func unmarshalEchoRequest(h header, buf []byte) (*EchoRequest, error) {
	if h.messageType != echoRequest {
		log.Fatal("Invalud messageType")
	}

	anIe, _, err := ie.Unmarshal(buf)
	if err != nil {
		return nil, err
	}
	rec := anIe.(*ie.Recovery)
	return NewEchoRequest(h.seqNum, rec.Value), nil
}

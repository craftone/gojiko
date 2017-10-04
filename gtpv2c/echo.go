package gtpv2c

import (
	"log"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type echo struct {
	header
	Recovery *ie.Recovery
}

type EchoRequest struct {
	echo
}

type EchoResponse struct {
	echo
}

func newEcho(messageType messageTypeNum, seqNum uint32, recoveryValue byte) (echo, error) {
	rec, err := ie.NewRecovery(0, recoveryValue)
	if err != nil {
		return echo{}, err
	}
	return echo{
		newHeader(messageType, false, false, 0, seqNum),
		rec,
	}, nil
}

func NewEchoRequest(seqNum uint32, recoveryValue byte) (*EchoRequest, error) {
	echo, err := newEcho(echoRequestNum, seqNum, recoveryValue)
	if err != nil {
		return nil, err
	}
	return &EchoRequest{echo}, nil
}

func NewEchoResponse(seqNum uint32, recoveryValue byte) (*EchoResponse, error) {
	echo, err := newEcho(echoResponseNum, seqNum, recoveryValue)
	if err != nil {
		return nil, err
	}
	return &EchoResponse{echo}, nil
}

func (e echo) Marshal() []byte {
	body := e.Recovery.Marshal()
	return e.header.marshal(body)
}

func unmarshalEchoRequest(h header, buf []byte) (*EchoRequest, error) {
	if h.messageType != echoRequestNum {
		log.Fatal("Invalud messageType")
	}

	anIe, _, err := ie.Unmarshal(buf, ie.EchoRequest)
	if err != nil {
		return nil, err
	}
	rec := anIe.(*ie.Recovery)
	return NewEchoRequest(h.seqNum, rec.Value)
}

func unmarshalEchoResponse(h header, buf []byte) (*EchoResponse, error) {
	if h.messageType != echoResponseNum {
		log.Fatal("Invalud messageType")
	}

	anIe, _, err := ie.Unmarshal(buf, ie.EchoResponse)
	if err != nil {
		return nil, err
	}
	rec := anIe.(*ie.Recovery)
	return NewEchoResponse(h.seqNum, rec.Value)
}

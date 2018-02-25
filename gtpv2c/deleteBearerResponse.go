package gtpv2c

import (
	"fmt"
	"log"

	"github.com/craftone/gojiko/gtp"
	"github.com/craftone/gojiko/gtpv2c/ie"
)

type DeleteBearerResponse struct {
	header
	cause *ie.Cause
}

func NewDeleteBearerResponse(teid gtp.Teid, seqNum uint32, cause ie.CauseValue) (*DeleteBearerResponse, error) {
	causeIE, err := ie.NewCause(0, cause, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &DeleteBearerResponse{
		newHeader(DeleteBearerResponseNum, false, true, teid, seqNum),
		causeIE,
	}, nil
}

func (d *DeleteBearerResponse) Marshal() []byte {
	return d.header.marshal(d.cause.Marshal())
}

func unmarshalDeleteBearerResponse(h header, buf []byte) (*DeleteBearerResponse, error) {
	if h.messageType != DeleteBearerResponseNum {
		panic(fmt.Sprintf("Invalid message Type : %d", h.messageType))
	}

	var causeIE *ie.Cause
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.DeleteBearerResponse)
		buf = tail
		if err != nil {
			if _, ok := err.(*ie.UnknownIEError); ok {
				log.Println(err)
				continue
			}
			return nil, err
		}

		if msg.Instance() != 0 {
			log.Printf("Unkown IE : %#v", msg)
		}

		switch msg := msg.(type) {
		case *ie.Cause:
			causeIE = msg
		default:
			log.Printf("Unkown IE : %#v", msg)
		}
	}
	if causeIE == nil {
		return nil, fmt.Errorf("No Cause Delete Bearer Request message")
	}
	return NewDeleteBearerResponse(h.Teid(), h.seqNum, causeIE.Value())
}

// Getters

func (d *DeleteBearerResponse) Cause() *ie.Cause {
	return d.cause
}

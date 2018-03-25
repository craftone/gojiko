package gtpv2c

import (
	"fmt"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type DeleteSessionResponse struct {
	header
	cause *ie.Cause
}

func NewDeleteSessionResponse(teid gtp.Teid, seqNum uint32, cause ie.CauseValue) (*DeleteSessionResponse, error) {
	causeIE, err := ie.NewCause(0, cause, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &DeleteSessionResponse{
		newHeader(DeleteSessionResponseNum, false, true, teid, seqNum),
		causeIE,
	}, nil
}

func (d *DeleteSessionResponse) Marshal() []byte {
	body := make([]byte, 0, 16)
	body = append(body, d.cause.Marshal()...)
	return d.header.marshal(body)
}

func unmarshalDeleteSessionResponse(h header, buf []byte) (*DeleteSessionResponse, error) {
	if h.messageType != DeleteSessionResponseNum {
		panic(fmt.Sprintf("Invalid message Type : %d", h.messageType))
	}

	var causeIE *ie.Cause
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.DeleteSessionResponse)
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
	return NewDeleteSessionResponse(h.Teid(), h.seqNum, causeIE.Value())
}

// Getters

func (d *DeleteSessionResponse) Cause() *ie.Cause {
	return d.cause
}

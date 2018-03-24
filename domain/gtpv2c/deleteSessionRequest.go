package gtpv2c

import (
	"fmt"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type DeleteSessionRequest struct {
	header
	lbi *ie.Ebi
}

func NewDeleteSessionRequest(teid gtp.Teid, seqNum uint32, ebi byte) (*DeleteSessionRequest, error) {
	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return nil, err
	}
	return &DeleteSessionRequest{
		newHeader(DeleteSessionRequestNum, false, true, teid, seqNum),
		ebiIE,
	}, nil
}

func (d *DeleteSessionRequest) Marshal() []byte {
	return d.header.marshal(d.lbi.Marshal())
}

func unmarshalDeleteSessionRequest(h header, buf []byte) (*DeleteSessionRequest, error) {
	if h.messageType != DeleteSessionRequestNum {
		panic(fmt.Sprintf("Invalid message Type : %d", h.messageType))
	}

	var ebiIE *ie.Ebi
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.DeleteSessionRequest)
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
		case *ie.Ebi:
			ebiIE = msg
		default:
			log.Printf("Unkown IE : %#v", msg)
		}
	}
	if ebiIE == nil {
		return nil, fmt.Errorf("No LBI (Linked EPS Bearer ID) Delete Bearer Request message")
	}
	return NewDeleteSessionRequest(h.Teid(), h.seqNum, ebiIE.Value())
}

// Getters

func (d *DeleteSessionRequest) Lbi() *ie.Ebi {
	return d.lbi
}

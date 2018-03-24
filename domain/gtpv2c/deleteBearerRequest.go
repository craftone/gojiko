package gtpv2c

import (
	"fmt"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type DeleteBearerRequest struct {
	header
	ebi *ie.Ebi
}

func NewDeleteBearerRequest(teid gtp.Teid, seqNum uint32, ebi byte) (*DeleteBearerRequest, error) {
	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return nil, err
	}
	return &DeleteBearerRequest{
		newHeader(DeleteBearerRequestNum, false, true, teid, seqNum),
		ebiIE,
	}, nil
}

func (d *DeleteBearerRequest) Marshal() []byte {
	return d.header.marshal(d.ebi.Marshal())
}

func unmarshalDeleteBearerRequest(h header, buf []byte) (*DeleteBearerRequest, error) {
	if h.messageType != DeleteBearerRequestNum {
		panic(fmt.Sprintf("Invalid message Type : %d", h.messageType))
	}

	var ebiIE *ie.Ebi
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.DeleteBearerRequest)
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
	return NewDeleteBearerRequest(h.Teid(), h.seqNum, ebiIE.Value())
}

// Getters

func (d *DeleteBearerRequest) Lbi() *ie.Ebi {
	return d.ebi
}

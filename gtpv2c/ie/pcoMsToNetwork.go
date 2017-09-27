package ie

import (
	"log"

	"github.com/craftone/gojiko/gtpv2c/ie/pco"
)

type PcoMsToNetwork struct {
	header
	*pco.MsToNetwork
}

func NewPcoMsToNetwork(instance byte, msToNetwork *pco.MsToNetwork) (*PcoMsToNetwork, error) {
	header, err := newHeader(pcoNum, 0, instance)
	if err != nil {
		return nil, err
	}
	return &PcoMsToNetwork{
		header,
		msToNetwork,
	}, nil
}

func (p *PcoMsToNetwork) Marshal() []byte {
	body := p.MsToNetwork.Marshal()
	return p.header.marshal(body)
}

func unmarshalPcoMsToNetwork(h header, buf []byte) (*PcoMsToNetwork, error) {
	if h.typeNum != pcoNum {
		log.Fatal("Invalud type")
	}

	msToNetwork, _, err := pco.UnmarshalMsToNetowrk(buf)
	if err != nil {
		return nil, err
	}

	pcoMsToNetwork, err := NewPcoMsToNetwork(h.instance, msToNetwork)
	if err != nil {
		return nil, err
	}
	return pcoMsToNetwork, nil
}

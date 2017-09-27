package ie

import (
	"log"

	"github.com/craftone/gojiko/gtpv2c/ie/pco"
)

type PcoNetworkToMs struct {
	header
	*pco.NetworkToMs
}

func NewPcoNetworkToMs(instance byte, networkToMs *pco.NetworkToMs) (*PcoNetworkToMs, error) {
	header, err := newHeader(pcoNum, 0, instance)
	if err != nil {
		return nil, err
	}
	return &PcoNetworkToMs{
		header,
		networkToMs,
	}, nil
}

func (p *PcoNetworkToMs) Marshal() []byte {
	body := p.NetworkToMs.Marshal()
	return p.header.marshal(body)
}

func unmarshalPcoNetworkToMs(h header, buf []byte) (*PcoNetworkToMs, error) {
	if h.typeNum != pcoNum {
		log.Fatal("Invalud type")
	}

	networkToMs, _, err := pco.UnmarshalNetowrkToMs(buf)
	if err != nil {
		return nil, err
	}

	pcoNetworkToMs, err := NewPcoNetworkToMs(h.instance, networkToMs)
	if err != nil {
		return nil, err
	}
	return pcoNetworkToMs, nil
}

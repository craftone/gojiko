package ie

import "log"

type ServingNetwork struct {
	header
	mccMnc
}

func NewServingNetwork(instance byte, mcc, mnc string) (*ServingNetwork, error) {
	header, err := newHeader(servingNetworkNum, 3, instance)
	if err != nil {
		return nil, err
	}
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &ServingNetwork{
		header,
		mccMnc,
	}, nil
}

func (r *ServingNetwork) Marshal() []byte {
	body := r.mccMnc.tbcd[:]
	return r.header.marshal(body)
}

func unmarshalServingNetwork(h header, buf []byte) (*ServingNetwork, error) {
	if h.typeNum != servingNetworkNum {
		log.Fatal("Invalud type")
	}

	mccMnc, _, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, err
	}
	sn, err := NewServingNetwork(h.instance, mccMnc.Mcc, mccMnc.Mnc)
	if err != nil {
		return nil, err
	}
	return sn, nil
}

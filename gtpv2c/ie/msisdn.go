package ie

import (
	"log"
)

type Msisdn struct {
	tbcdIE
}

const MSISDN_MAX_LEN = 15

func NewMsisdn(instance byte, value string) (*Msisdn, error) {
	tbcdIE, err := newTbcdIE(msisdnNum, 0, instance, value, 6, 15)
	if err != nil {
		return nil, err
	}
	return &Msisdn{tbcdIE}, nil
}

func (m *Msisdn) Marshal() []byte {
	return m.tbcdIE.marshal()
}

func unmarshalMsisdn(h header, buf []byte) (*Msisdn, error) {
	if h.typeNum != msisdnNum {
		log.Fatal("Invalud type")
	}

	s, err := unmarshalTbcd(buf)
	if err != nil {
		return nil, err
	}
	msisdn, err := NewMsisdn(h.instance, s)
	if err != nil {
		return nil, err
	}
	return msisdn, nil
}

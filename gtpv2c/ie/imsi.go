package ie

import (
	"log"
)

type Imsi struct {
	header
	Value string
	tbcd  tbcd
}

func NewImsi(instance byte, value string) (*Imsi, error) {
	if len(value) < 6 || len(value) > 15 {
		log.Fatal("Number of IMSI digits must be from 6 to 15")
	}
	tbcd, err := parseTBCD(value)
	if err != nil {
		log.Fatal("Invalid imsi")
	}

	header, err := newHeader(imsiNum, uint16(len(value)), instance)
	if err != nil {
		return nil, err
	}

	return &Imsi{
		header: header,
		Value:  value,
		tbcd:   tbcd,
	}, nil
}

func (i *Imsi) Marshal() []byte {
	return i.header.marshal(i.tbcd)
}

func unmarshalImsi(h header, buf []byte) (*Imsi, error) {
	if h.typeNum != imsiNum {
		log.Fatal("Invalud type")
	}

	s, err := unmarshalTbcd(buf)
	if err != nil {
		return nil, err
	}
	imsi, err := NewImsi(h.instance, s)
	if err != nil {
		return nil, err
	}
	return imsi, nil
}

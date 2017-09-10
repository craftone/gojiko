package ie

import (
	"errors"
	"log"
)

type IndicationArg struct {
	DAF   bool
	DTF   bool
	HI    bool
	DFI   bool
	OI    bool
	ISRSI bool
	ISRAI bool
	SGWCI bool
	SQCI  bool
	UIMSI bool
	CFSI  bool
	CRSI  bool
	P     bool
	PT    bool
	SI    bool
	MSV   bool
	ISRAU bool
	CCRSI bool
}

type Indication struct {
	header
	IndicationArg
}

func NewIndication(instance byte, IndicationArg IndicationArg) (*Indication, error) {
	header, err := newHeader(indicationNum, 3, instance)
	if err != nil {
		return nil, err
	}
	return &Indication{header, IndicationArg}, nil
}

func (i *Indication) Marshal() []byte {
	body := make([]byte, 3)
	body[0] = setBit(body[0], 7, i.DAF)
	body[0] = setBit(body[0], 6, i.DTF)
	body[0] = setBit(body[0], 5, i.HI)
	body[0] = setBit(body[0], 4, i.DFI)
	body[0] = setBit(body[0], 3, i.OI)
	body[0] = setBit(body[0], 2, i.ISRSI)
	body[0] = setBit(body[0], 1, i.ISRAI)
	body[0] = setBit(body[0], 0, i.SGWCI)
	body[1] = setBit(body[1], 7, i.SQCI)
	body[1] = setBit(body[1], 6, i.UIMSI)
	body[1] = setBit(body[1], 5, i.CFSI)
	body[1] = setBit(body[1], 4, i.CRSI)
	body[1] = setBit(body[1], 3, i.P)
	body[1] = setBit(body[1], 2, i.PT)
	body[1] = setBit(body[1], 1, i.SI)
	body[1] = setBit(body[1], 0, i.MSV)
	body[2] = setBit(body[2], 1, i.ISRAU)
	body[2] = setBit(body[2], 0, i.CCRSI)

	return i.header.marshal(body)
}

func unmarshalIndication(h header, buf []byte) (*Indication, error) {
	if h.typeNum != indicationNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 3 {
		return nil, errors.New("Invalid binary")
	}

	indicationArg := IndicationArg{}

	indicationArg.DAF = getBit(buf[0], 7)
	indicationArg.DTF = getBit(buf[0], 6)
	indicationArg.HI = getBit(buf[0], 5)
	indicationArg.DFI = getBit(buf[0], 4)
	indicationArg.OI = getBit(buf[0], 3)
	indicationArg.ISRSI = getBit(buf[0], 2)
	indicationArg.ISRAI = getBit(buf[0], 1)
	indicationArg.SGWCI = getBit(buf[0], 0)
	indicationArg.SQCI = getBit(buf[1], 7)
	indicationArg.UIMSI = getBit(buf[1], 6)
	indicationArg.CFSI = getBit(buf[1], 5)
	indicationArg.CRSI = getBit(buf[1], 4)
	indicationArg.P = getBit(buf[1], 3)
	indicationArg.PT = getBit(buf[1], 2)
	indicationArg.SI = getBit(buf[1], 1)
	indicationArg.MSV = getBit(buf[1], 0)
	indicationArg.ISRAU = getBit(buf[2], 1)
	indicationArg.CCRSI = getBit(buf[2], 0)
	return NewIndication(h.instance, indicationArg)
}

package ie

import (
	"errors"
	"log"
)

type Cause struct {
	header
	Value       byte
	Pce         bool
	Bce         bool
	Cs          bool
	OffendingIe *header
}

func NewCause(instance byte, value byte, pce, bce, cs bool, offendingIe *header) (*Cause, error) {
	length := 2
	if offendingIe != nil {
		length = 6
	}

	header, err := newHeader(causeNum, uint16(length), instance)
	if err != nil {
		return nil, err
	}

	return &Cause{
		header:      header,
		Value:       value,
		Pce:         pce,
		Bce:         bce,
		Cs:          cs,
		OffendingIe: offendingIe,
	}, nil
}

func (c *Cause) Marshal() []byte {
	buf := make([]byte, c.header.length)
	buf[0] = c.Value
	buf[1] = setBit(buf[1], 2, c.Pce)
	buf[1] = setBit(buf[1], 1, c.Bce)
	buf[1] = setBit(buf[1], 0, c.Cs)
	if c.OffendingIe != nil {
		buf[2] = byte(c.OffendingIe.typeNum)
		buf[5] = c.OffendingIe.instance
	}
	return c.header.marshal(buf)
}

func unmarshalCause(h header, buf []byte) (*Cause, error) {
	if h.typeNum != causeNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 2 {
		return nil, errors.New("too short data")
	}
	value := buf[0]
	pce := getBit(buf[1], 2)
	bce := getBit(buf[1], 1)
	cs := getBit(buf[1], 0)

	var offendingIeHeader *header
	if h.length == 6 {
		offendingIeHeader = &header{
			typeNum:  ieTypeNum(buf[2]),
			length:   0,
			instance: buf[5],
		}
	}
	cause, err := NewCause(h.instance, value, pce, bce, cs, offendingIeHeader)
	if err != nil {
		return nil, err
	}
	return cause, nil
}

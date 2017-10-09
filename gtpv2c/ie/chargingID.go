package ie

import (
	"encoding/binary"
	"errors"
	"log"
)

type ChargingID struct {
	header
	value uint32
}

func NewChargingID(instance byte, value uint32) (*ChargingID, error) {
	header, err := newHeader(chargingIDNum, 4, instance)
	if err != nil {
		return nil, err
	}

	return &ChargingID{
		header: header,
		value:  value,
	}, nil
}

func (c *ChargingID) Marshal() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf[0:4], c.value)
	return c.header.marshal(buf)
}

func unmarshalChargingID(h header, buf []byte) (*ChargingID, error) {
	if h.typeNum != chargingIDNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 4 {
		return nil, errors.New("too short data")
	}
	value := binary.BigEndian.Uint32(buf[0:4])

	chargingID, err := NewChargingID(h.instance, value)
	if err != nil {
		return nil, err
	}
	return chargingID, nil
}

func (c *ChargingID) Value() uint32 {
	return c.value
}

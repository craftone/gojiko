package ie

import (
	"errors"
	"fmt"
)

type RatType struct {
	header
	value byte
}

func NewRatType(instance byte, value byte) (*RatType, error) {
	header, err := newHeader(ratTypeNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &RatType{
		header: header,
		value:  value,
	}, nil
}

func (r *RatType) Marshal() []byte {
	body := []byte{r.value}
	return r.header.marshal(body)
}

func unmarshalRatType(h header, buf []byte) (*RatType, error) {
	if h.typeNum != ratTypeNum {
		log.Panic("Invalid type")
	}

	if len(buf) < 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewRatType(h.instance, buf[0])
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *RatType) Value() byte {
	return r.value
}

func (r *RatType) String() string {
	var str string
	switch int(r.value) {
	case 1:
		str = "UTRAN"
	case 2:
		str = "GERAN"
	case 3:
		str = "WLAN"
	case 4:
		str = "GAN"
	case 5:
		str = "HSPA Evolution"
	case 6:
		str = "EUTRAN"
	case 7:
		str = "Virtual"
	default:
		str = "<reserved>"
	}
	return fmt.Sprintf("%s (%d)", str, r.value)
}

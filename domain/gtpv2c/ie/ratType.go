package ie

import (
	"errors"
	"fmt"
)

type RatType struct {
	header
	value RatTypeValue
}

type RatTypeValue byte

const (
	RatTypeUtran         RatTypeValue = 1
	RatTypeGeran         RatTypeValue = 2
	RatTypeWlan          RatTypeValue = 3
	RatTypeGan           RatTypeValue = 4
	RatTypeHspaEvolution RatTypeValue = 5
	RatTypeWbEutran      RatTypeValue = 6
	RatTypeVirtual       RatTypeValue = 7
	RatTypeEutranNBIoT   RatTypeValue = 8
	RatTypeLteM          RatTypeValue = 9
	RatTypeNR            RatTypeValue = 10
)

func NewRatType(instance byte, value RatTypeValue) (*RatType, error) {
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
	body := []byte{byte(r.value)}
	return r.header.marshal(body)
}

func unmarshalRatType(h header, buf []byte) (*RatType, error) {
	if h.typeNum != ratTypeNum {
		log.Panic("Invalid type")
	}

	if len(buf) < 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewRatType(h.instance, RatTypeValue(buf[0]))
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *RatType) Value() RatTypeValue {
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
		str = "EUTRAN (WB-E-UTRAN)"
	case 7:
		str = "Virtual"
	case 8:
		str = "EUTRAN-NB-IoT"
	case 9:
		str = "LTE-M"
	case 10:
		str = "NR"
	default:
		str = "<reserved>"
	}
	return fmt.Sprintf("%s (%d)", str, r.value)
}

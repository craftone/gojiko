package ie

import "log"
import "errors"

type RatType struct {
	header
	Value byte
}

func NewRatType(instance byte, value byte) (*RatType, error) {
	header, err := newHeader(ratTypeNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &RatType{
		header: header,
		Value:  value,
	}, nil
}

func (r *RatType) Marshal() []byte {
	body := []byte{r.Value}
	return r.header.marshal(body)
}

func unmarshalRatType(h header, buf []byte) (*RatType, error) {
	if h.typeNum != ratTypeNum {
		log.Fatal("Invalud type")
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

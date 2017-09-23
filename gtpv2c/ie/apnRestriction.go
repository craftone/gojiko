package ie

import "log"
import "errors"

type ApnRestriction struct {
	header
	Value byte
}

func NewApnRestriction(instance byte, value byte) (*ApnRestriction, error) {
	header, err := newHeader(apnRestrictionNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &ApnRestriction{
		header: header,
		Value:  value,
	}, nil
}

func (r *ApnRestriction) Marshal() []byte {
	body := []byte{r.Value}
	return r.header.marshal(body)
}

func unmarshalApnRestriction(h header, buf []byte) (*ApnRestriction, error) {
	if h.typeNum != apnRestrictionNum {
		log.Fatal("Invalud type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewApnRestriction(h.instance, buf[0])
	if err != nil {
		return nil, err
	}
	return rec, nil
}

package ie

import "errors"

type Recovery struct {
	header
	value byte
}

func NewRecovery(instance byte, value byte) (*Recovery, error) {
	header, err := newHeader(recoveryNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &Recovery{
		header: header,
		value:  value,
	}, nil
}

func (r *Recovery) Marshal() []byte {
	body := []byte{r.value}
	return r.header.marshal(body)
}

func unmarshalRecovery(h header, buf []byte) (*Recovery, error) {
	if h.typeNum != recoveryNum {
		log.Panic("Invalid type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewRecovery(h.instance, buf[0])
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *Recovery) Value() byte {
	return r.value
}

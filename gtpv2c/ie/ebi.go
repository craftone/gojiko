package ie

import "log"
import "errors"
import "fmt"

type Ebi struct {
	header
	value byte
}

func NewEbi(instance byte, value byte) (*Ebi, error) {
	if value < 5 || value > 15 {
		return nil, fmt.Errorf("Unkown EBI : %d", value)
	}
	header, err := newHeader(ebiNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &Ebi{
		header: header,
		value:  value,
	}, nil
}

func (r *Ebi) Marshal() []byte {
	body := []byte{r.value}
	return r.header.marshal(body)
}

func unmarshalEbi(h header, buf []byte) (*Ebi, error) {
	if h.typeNum != ebiNum {
		log.Fatal("Invalud type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewEbi(h.instance, buf[0])
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (e *Ebi) Value() byte {
	return e.value
}

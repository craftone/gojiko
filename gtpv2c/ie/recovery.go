package ie

import "log"
import "errors"

type Recovery struct {
	*header
	Value byte
}

func NewRecovery(instance byte, value byte) *Recovery {
	return &Recovery{
		header: newHeader(recoveryNum, 5, instance),
		Value:  value,
	}
}

func (r *Recovery) Marshal() []byte {
	body := []byte{r.Value}
	return r.header.marshal(body)
}

func unmarshalRecovery(h header, buf []byte) (*Recovery, error) {
	if h.typeNum != recoveryNum {
		log.Fatal("Invalud type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}
	return NewRecovery(h.instance, buf[0]), nil
}

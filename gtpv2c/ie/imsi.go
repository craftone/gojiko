package ie

import (
	"log"
)

type Imsi struct {
	header *header
	Value  string
}

func NewImsi(value string, instance byte) *Imsi {
	if len(value) < 6 || len(value) > 15 {
		log.Fatal("Number of IMSI digits must be from 6 to 16")
	}
	return &Imsi{
		header: newHeader(imsiNum, uint16(len(value)), instance),
		Value:  value,
	}
}

// func (r *Imsi) Marshal() []byte {
// 	body := []byte{r.Value}
// 	return r.header.marshal(body)
// }

// func unmarshalImsi(h header, buf []byte) (*Imsi, error) {
// 	if h.typeNum != ImsiNum {
// 		log.Fatal("Invalud type")
// 	}

// 	if len(buf) != 1 {
// 		return nil, errors.New("Invalid binary")
// 	}
// 	return NewImsi(buf[0], h.instance), nil
// }

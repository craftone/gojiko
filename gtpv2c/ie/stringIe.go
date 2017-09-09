package ie

import (
	"log"
)

type stringIE struct {
	header
	Value string
}

type Apn struct {
	stringIE
}

func newStringIE(typeNum ieTypeNum, length uint16, instance byte, value string, minLen, maxLen int) stringIE {
	if len(value) < minLen || len(value) > maxLen {
		log.Fatalf("Number of digits must be from %d to %d", minLen, maxLen)
	}
	return stringIE{
		header: newHeader(typeNum, length, instance),
		Value:  value,
	}
}

func (s *stringIE) marshal() []byte {
	return s.header.marshal([]byte(s.Value))
}

func NewApn(instance byte, value string) *Apn {
	return &Apn{
		stringIE: newStringIE(apnNum, 0, instance, value, 1, 100),
	}
}

func (i *Apn) Marshal() []byte {
	return i.stringIE.marshal()
}

func unmarshalApn(h header, buf []byte) (*Apn, error) {
	if h.typeNum != apnNum {
		log.Fatal("Invalud type")
	}

	return NewApn(h.instance, string(buf)), nil
}

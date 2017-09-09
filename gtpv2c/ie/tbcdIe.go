package ie

import (
	"fmt"
)

type tbcdIE struct {
	header
	Value string
	tbcd  tbcd
}

func newTbcdIE(typeNum ieTypeNum, length uint16, instance byte, value string, minLen, maxLen int) (tbcdIE, error) {
	if len(value) < minLen || len(value) > maxLen {
		return tbcdIE{}, fmt.Errorf("Number of digits must be from %d to %d", minLen, maxLen)
	}
	tbcd, err := parseTBCD(value)
	if err != nil {
		return tbcdIE{}, err
	}

	header, err := newHeader(typeNum, length, instance)
	if err != nil {
		return tbcdIE{}, err
	}

	return tbcdIE{
		header: header,
		Value:  value,
		tbcd:   tbcd,
	}, nil
}

func (t *tbcdIE) marshal() []byte {
	return t.header.marshal(t.tbcd)
}

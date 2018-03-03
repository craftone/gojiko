package ie

import (
	"fmt"
)

type tbcdIE struct {
	header
	value string
	tbcd  tbcd
}

func newTbcdIE(typeNum ieTypeNum, addLen uint16, instance byte, value string, minLen, maxLen int) (tbcdIE, error) {
	if len(value) < minLen || len(value) > maxLen {
		return tbcdIE{}, fmt.Errorf("Number of digits must be from %d to %d", minLen, maxLen)
	}
	tbcd, err := parseTBCD(value)
	if err != nil {
		return tbcdIE{}, err
	}

	header, err := newHeader(typeNum, uint16(len(tbcd))+addLen, instance)
	if err != nil {
		return tbcdIE{}, err
	}

	return tbcdIE{
		header: header,
		value:  value,
		tbcd:   tbcd,
	}, nil
}

func (t *tbcdIE) marshal() []byte {
	return t.header.marshal(t.tbcd)
}

func (t *tbcdIE) Value() string {
	return t.value
}

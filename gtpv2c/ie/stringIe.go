package ie

import (
	"fmt"
)

type stringIE struct {
	header
	Value string
}

func newStringIE(typeNum ieTypeNum, length uint16, instance byte, value string, minLen, maxLen int) (stringIE, error) {
	if len(value) < minLen || len(value) > maxLen {
		return stringIE{}, fmt.Errorf("Number of digits must be from %d to %d", minLen, maxLen)
	}

	header, err := newHeader(typeNum, length, instance)
	if err != nil {
		return stringIE{}, err
	}

	return stringIE{
		header: header,
		Value:  value,
	}, nil
}

func (s *stringIE) marshal() []byte {
	return s.header.marshal([]byte(s.Value))
}
